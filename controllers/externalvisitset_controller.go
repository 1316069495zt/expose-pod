/*
Copyright 2021.

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
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "external/api/v1alpha1"
)

// ExternalvisitSetReconciler reconciles a ExternalvisitSet object
type ExternalvisitSetReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	processor *Processor
}

//+kubebuilder:rbac:groups=apps.zt1,resources=externalvisitsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.zt1,resources=externalvisitsets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps.zt1,resources=externalvisitsets/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ExternalvisitSet object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *ExternalvisitSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	externalvisitSet := &appsv1alpha1.ExternalvisitSet{}

	err := r.Get(context.TODO(), req.NamespacedName, externalvisitSet)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	// your logic here
	r.processor = NewSidecarSetProcessor(r.Client, nil)
	return r.processor.UpdateExternalvisitSet(externalvisitSet)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExternalvisitSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.ExternalvisitSet{}).
		Complete(r)
}
