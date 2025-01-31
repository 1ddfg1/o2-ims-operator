/*
Copyright 2025.

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
	"io"
	"net/http"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	ovranv1alpha1 "github.com/example/o2-ims-operator/api/v1alpha1"
)

const (
	typeAvailableO2ims = "Available"
	typeDegradedO2ims  = "Degraded"
)

// O2imsReconciler reconciles a O2ims object
type O2imsReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=ovran.example.com,resources=o2ims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ovran.example.com,resources=o2ims/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ovran.example.com,resources=o2ims/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *O2imsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	o2ims := &ovranv1alpha1.O2ims{}
	err := r.Get(ctx, req.NamespacedName, o2ims)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("o2ims resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get o2ims")
		return ctrl.Result{}, err
	}

	if len(o2ims.Status.Conditions) == 0 {
		meta.SetStatusCondition(&o2ims.Status.Conditions, metav1.Condition{Type: typeAvailableO2ims, Status: metav1.ConditionUnknown, Reason: "Reconciling", Message: "Starting reconciliation"})
		if err = r.Status().Update(ctx, o2ims); err != nil {
			log.Error(err, "Failed to update O2ims status")
			return ctrl.Result{}, err
		}
		if err := r.Get(ctx, req.NamespacedName, o2ims); err != nil {
			log.Error(err, "Failed to re-fetch O2ims")
			return ctrl.Result{}, err
		}
	}

	isNginxMarkedToBeDeleted := o2ims.GetDeletionTimestamp() != nil
	if isNginxMarkedToBeDeleted {
		return ctrl.Result{}, nil
	}

	log.Info("URL is : " + o2ims.Spec.Url)
	response, err := http.Get(o2ims.Spec.Url)
	if err != nil {
		log.Error(err, "Failed to invoke REST endpoint: "+o2ims.Spec.Url)
		return ctrl.Result{}, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error(err, "Failed to read REST response from URL: "+o2ims.Spec.Url)
		return ctrl.Result{}, err
	}

	log.Info("Successfully invoked URL: " + o2ims.Spec.Url)
	log.Info(string(responseData))

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *O2imsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ovranv1alpha1.O2ims{}).
		Complete(r)
}
