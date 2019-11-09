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

package controllers

import (
	"context"
	"sync"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/KeisukeYamashita/i/api/v1alpha1"
	icontrollerv1alpha1 "github.com/KeisukeYamashita/i/api/v1alpha1"
	"github.com/KeisukeYamashita/i/pkgs/syncers"
)

// EyeReconciler reconciles a Eye object
type EyeReconciler struct {
	client.Client
	Log      logr.Logger
	Schema   *runtime.Scheme
	Recorder record.EventRecorder

	syncers map[types.NamespacedName]syncers.Syncer
	mu      sync.RWMutex
}

// +kubebuilder:rbac:groups=icontroller.i.keisukeyamashita.com,resources=eyes,verbs=get;list;watch;create;update;patchwatch;list
// +kubebuilder:rbac:groups=icontroller.i.keisukeyamashita.com,resources=eyes/status,verbs=get;update;patch

func (r *EyeReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	nn := req.NamespacedName
	ctx := context.Background()
	log := r.Log.WithValues("eye", nn)

	var eye v1alpha1.Eye
	r.Log.Info("fetching eye object")
	if err := r.Get(ctx, req.NamespacedName, &eye); err != nil {
		log.Error(err, "unable to fetch eye")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	r.mu.RLock()
	_, ok := r.syncers[nn]
	r.mu.RUnlock()
	if !ok {

	}

	if eye.Spec.ExpireAt == "" {
		log.V(1).Info("not eye target")
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

func (r *EyeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&icontrollerv1alpha1.Eye{}).
		Complete(r)
}
