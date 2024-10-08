package syncer

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/KeisukeYamashita/i/api/v1alpha1"
	"github.com/KeisukeYamashita/i/pkg/clock"
	"github.com/KeisukeYamashita/i/pkg/logging"
	"github.com/KeisukeYamashita/i/pkg/slack"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

// Syncer ...
type Syncer interface {
	Sync(context.Context)
	Stop()
}

type syncer struct {
	ticker     *time.Ticker
	client     client.Client
	lifetime   time.Duration
	clock      clock.Clock
	HookURL    *url.URL
	eye        *v1alpha1.Eye
	nn         types.NamespacedName
	CancelFunc context.CancelFunc
}

var _ Syncer = (*syncer)(nil)

// NewSyncer ...
func NewSyncer(
	client client.Client,
	eye *v1alpha1.Eye,
	nn types.NamespacedName,
	clock clock.Clock,
	url *url.URL,
	cancelFunc context.CancelFunc,
) (Syncer, error) {
	d, err := time.ParseDuration(eye.Spec.Lifetime)
	if err != nil {
		return nil, err
	}
	return &syncer{
		client:     client,
		clock:      clock,
		lifetime:   d,
		eye:        eye,
		nn:         nn,
		HookURL:    url,
		CancelFunc: cancelFunc,
	}, nil
}

// Sync ...
func (s *syncer) Sync(ctx context.Context) {
	log := logging.FromContext(ctx)
	ticker := time.NewTicker(time.Duration(1 * time.Second))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.V(0).Info("trigger check", "name", s.eye.Name, "lifetime", s.eye.Spec.Lifetime)
			eg, _ := errgroup.WithContext(ctx)
			pods := []corev1.Pod{}
			eg.Go(func() error {
				var err error
				if pods, err = s.getPods(ctx); err != nil {
					return err
				}
				return nil
			})

			if err := eg.Wait(); err != nil {
				return
			}

			invalidPods := []corev1.Pod{}
			for _, pod := range pods {
				if valid := s.validPod(&pod); valid {
					pod := pod
					invalidPods = append(invalidPods, pod)
					eg.Go(func() error {
						err := s.deletePod(ctx, &pod)
						if err != nil {
							log.Error(err, fmt.Sprintf("deleted pod %s", pod.Name))
							return err
						}
						log.V(0).Info("delete pod", "pod", pod.Name)
						return nil
					})
				}
			}

			if err := eg.Wait(); err != nil {
				return
			}

			if len(invalidPods) != 0 {
				if s.HookURL != nil {
					msg := slack.NewInvalidPodsMessage(s.eye, s.nn, invalidPods)
					slackClient := slack.NewClient(s.HookURL)
					slackClient.PostMessage(msg)
				}
			}
		case <-ctx.Done():
			log.V(1).Info("stopping syncer with cancel", "err", ctx.Err())
			return
		}
	}
}

func (s *syncer) getPods(ctx context.Context) ([]corev1.Pod, error) {
	pl := &corev1.PodList{}
	if err := s.client.List(ctx, pl); err != nil {
		return nil, err
	}

	return pl.Items, nil
}

// validPod ...
func (s *syncer) validPod(pod *corev1.Pod) bool {
	now := s.clock.Now()
	expiresAt := pod.ObjectMeta.CreationTimestamp.Add(s.lifetime)
	return now.After(expiresAt)
}

func (s *syncer) deletePod(ctx context.Context, pod *corev1.Pod) error {
	if err := s.client.Delete(ctx, pod); err != nil {
		return err
	}

	return nil
}

func (s *syncer) Stop() {
	s.CancelFunc()
	return
}
