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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var machinemetadatalog = logf.Log.WithName("machinemetadata-resource")

func (r *MachineMetadata) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-x-cluster-x-k8s-io-v1alpha1-machinemetadata,mutating=true,failurePolicy=fail,groups=x-cluster.x-k8s.io,resources=machinemetadata,verbs=create;update,versions=v1alpha1,name=mmachinemetadata.kb.io

var _ webhook.Defaulter = &MachineMetadata{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *MachineMetadata) Default() {
	machinemetadatalog.Info("default", "name", r.Name)

	if r.Spec.Lables == nil {
		r.Spec.Lables = map[string]string{}
	}
	if r.Spec.Annotations == nil {
		r.Spec.Annotations = map[string]string{}
	}
	if r.Spec.Taints == nil {
		r.Spec.Taints = []corev1.Taint{}
	}
	if r.Status.Targets == nil {
		r.Status.Targets = []string{}
	}
}
