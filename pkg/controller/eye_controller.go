/*
Copyright 2019 KeisukeYamashita.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"net/url"
	"sync"
	"time"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/KeisukeYamashita/i/api/v1alpha1"
	icontrollerv1alpha1 "github.com/KeisukeYamashita/i/api/v1alpha1"
	"github.com/KeisukeYamashita/i/pkg/clock"
	"github.com/KeisukeYamashita/i/pkg/logging"
	"github.com/KeisukeYamashita/i/pkg/syncer"
)

// EyeReconciler reconciles a Eye object
type EyeReconciler struct {
	client.Client
	Log      logr.Logger
	Schema   *runtime.Scheme
	Recorder record.EventRecorder
	HookURL  *url.URL

	syncers map[types.NamespacedName]syncer.Syncer
	mu      sync.RWMutex
	clock.Clock
}

// NewEyeReconciler ...
func NewEyeReconciler(client client.Client, logger logr.Logger) *EyeReconciler {
	clock := clock.NewClock(time.Now())
	return &EyeReconciler{
		Client:  client,
		Log:     logger,
		Clock:   clock,
		syncers: make(map[types.NamespacedName]syncers.Syncer),
	}
}

// +kubebuilder:rbac:groups=icontroller.i.keisukeyamashita.com,resources=eyes,verbs=get;list;watch;create;update;patchwatch;list
// +kubebuilder:rbac:groups=icontroller.i.keisukeyamashita.com,resources=eyes/status,verbs=get;update;patch

// Reconcile handles the control-loop for pods
func (r *EyeReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	nn := req.NamespacedName
	ctx := context.Background()
	log := r.Log.WithValues("eye", nn)

	r.mu.RLock()
	s, ok := r.syncers[nn]
	r.mu.RUnlock()

	var eye v1alpha1.Eye
	r.Log.Info("fetching eye object")
	if err := r.Get(ctx, req.NamespacedName, &eye); err != nil {
		log.Error(err, "unable to fetch eye")
		s.Stop()
		r.mu.Lock()
		delete(r.syncers, nn)
		r.mu.Unlock()
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	url, err := r.GetSecret(ctx, eye.Spec.SecretRef.Name, &nn)
	if err != nil {
		r.Log.Error(err, "secret not found")
	}
	if url != nil {
		r.HookURL = url
	}

	if !ok {
		ctx = logging.WithContext(ctx, log)
		err := r.startSyncer(ctx, r.Client, nn, &eye)
		return ctrl.Result{}, client.IgnoreNotFound((err))
	}

	log.V(0).Info("update eye resource")
	return ctrl.Result{}, nil
}

// SetupWithManager ...
func (r *EyeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&icontrollerv1alpha1.Eye{}).
		Complete(r)
}

func (r *EyeReconciler) startSyncer(ctx context.Context, c client.Client, nn types.NamespacedName, eye *v1alpha1.Eye) error {
	log := logging.FromContext(ctx)
	log.V(1).Info("adding new syncer")
	ctx, cancel := context.WithCancel(ctx)
	s, err := syncers.NewSyncer(c, eye, nn, r.Clock, r.HookURL, cancel)
	if err != nil {
		return err
	}
	go s.Sync(ctx)
	r.mu.Lock()
	r.syncers[nn] = s
	r.mu.Unlock()

	log.V(1).Info("add new syncer")
	return nil
}

// GetSecret ...
func (r *EyeReconciler) GetSecret(ctx context.Context, name string, nn *types.NamespacedName) (*url.URL, error) {
	secret := &corev1.Secret{}
	// Copy types.NamedSpaces
	nn2 := *nn
	nn2.Name = name
	if err := r.Client.Get(ctx, nn2, secret); err != nil {
		return nil, err
	}

	data, ok := secret.Data["SLACK_URL"]
	if !ok {
		return nil, nil
	}

	url, err := url.Parse(string(data))
	if err != nil {
		return nil, err
	}

	return url, nil
}
