package syncers

import (
	"context"
	"time"

	"github.com/KeisukeYamashita/i/api/v1alpha1"
	"github.com/KeisukeYamashita/i/pkgs/clock"
	"github.com/k0kubun/pp"
	"golang.org/x/sync/errgroup"
	corev1 "k8s.io/api/core/v1"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

// Syncer ...
type Syncer interface {
	Sync(context.Context)
}

type syncer struct {
	ticker     *time.Ticker
	client     client.Client
	lifetime   time.Duration
	clock      clock.Clock
	CancelFunc context.CancelFunc
}

var _ Syncer = (*syncer)(nil)

// NewSyncer ...
func NewSyncer(client client.Client, eye *v1alpha1.Eye, clock clock.Clock, cancelFunc context.CancelFunc) (Syncer, error) {
	d, err := time.ParseDuration(eye.Spec.Lifetime)
	if err != nil {
		return nil, err
	}
	return &syncer{
		client:     client,
		clock:      clock,
		lifetime:   d,
		CancelFunc: cancelFunc,
	}, nil
}

// Sync ...
func (s *syncer) Sync(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(3 * time.Second))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
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
					invalidPods = append(invalidPods, pod)
				}
			}

			pp.Println(len(invalidPods))

		case <-ctx.Done():
			pp.Println(ctx.Err())
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
	return expiresAt.After(now)
}
