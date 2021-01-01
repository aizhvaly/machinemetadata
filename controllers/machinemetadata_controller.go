/*

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

	xclusterv1 "github.com/aizhvaly/machinemetadata/api/v1alpha1"
	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MachineMetadataReconciler reconciles a MachineMetadata object
type MachineMetadataReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=x-cluster.x-k8s.io,resources=machinemetadata,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=x-cluster.x-k8s.io,resources=machinemetadata/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machinedeployments;machines,verbs=get;list;watch

func (r *MachineMetadataReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("machinemetadata", req.NamespacedName)

	var machineMetadata xclusterv1.MachineMetadata
	if err := r.Client.Get(ctx, req.NamespacedName, &machineMetadata); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("machinemetadata not found, nothing to reconcile", "key", req.NamespacedName)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	var machineList clusterv1.MachineList
	labels := map[string]string{clusterv1.MachineDeploymentLabelName: machineMetadata.Spec.MachineDeploymentName}
	if err := r.Client.List(ctx, &machineList, client.InNamespace(req.Namespace), client.MatchingLabels(labels)); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("machines in machinedeployment not found, nothing to reconcile", "key", machineMetadata.Spec.MachineDeploymentName)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	for _, ma := range machineList.Items {
		if ma.Status.NodeRef != nil {
			machineMetadata.Status.Targets = append(machineMetadata.Status.Targets, ma.Status.NodeRef.Name)
			log.Info("discovered machine", "machine", ma.Status.NodeRef.Name)
		}
	}

	if err := r.Status().Update(ctx, &machineMetadata); err != nil {
		log.Error(err, "unable to update machinemetadata status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *MachineMetadataReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&xclusterv1.MachineMetadata{}).
		Complete(r)
}
