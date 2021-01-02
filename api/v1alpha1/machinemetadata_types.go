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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MachineMetadataSpec defines the desired state of MachineMetadata.
type MachineMetadataSpec struct {
	// MachineDeploymentName is the name of the MachineDeployment this object belongs to.
	MachineDeploymentName string `json:"machineDeploymentName"`

	// Cluster name is the name of the cluster in which needed reconcile (also mmd will find kubeconfig def capi secrets for cluster "clustername-kubeconfig" as example).
	// +optional
	ClusterName string `json:"cluster,omitempty"`
	// Lables is the labels which must be setted on nodes of MachineDeployment.
	// +optional
	Lables map[string]string `json:"labels,omitempty"`

	// Annotations is the annotations which must be setted on nodes of MachineDeployment.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// Taints is the taints which must be setted on nodes of MachineDeployment.
	// +optional
	Taints []corev1.Taint `json:"taints,omitempty"`
}

// MachineMetadataStatus defines the observed state of MachineMetadata
type MachineMetadataStatus struct {
	// Targets shows the current list of machines the machine metadata is watching.
	// +optional
	Targets []string `json:"targets,omitempty"`

	// Status defines the observed state of MachineMetadata
	// +optional
	Status Status `json:"status,omitempty"`
}

type StatusResult string

const (
	ResultSuccess StatusResult = "Success"
	ResultFail    StatusResult = "Fail"
)

// Status defines the observed state of MachineMetadata
type Status struct {
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`
	LastResult         StatusResult `json:"lastResult,omitempty"`
	Msg                string       `json:"msg,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=mmd
// +kubebuilder:subresource:status

// MachineMetadata is the Schema for the machinemetadata API
type MachineMetadata struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MachineMetadataSpec   `json:"spec,omitempty"`
	Status MachineMetadataStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MachineMetadataList contains a list of MachineMetadata
type MachineMetadataList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MachineMetadata `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MachineMetadata{}, &MachineMetadataList{})
}
