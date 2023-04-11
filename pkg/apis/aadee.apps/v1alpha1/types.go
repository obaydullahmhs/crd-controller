// +groupName=aadee.apps
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ContainerSpec struct {
	Image string `json:"image,omitempty"`
	Port  int32  `json:"port,omitempty"`
}

// AadeeSpec defines the desired state of AadeeCRD
type AadeeSpec struct {
	Name      string        `json:"name,omitempty"`
	Replicas  *int32        `json:"replicas"`
	Container ContainerSpec `json:"container,container"`
}

// AadeeStatus defines the observed state of AadeeCRD
type AadeeStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// Aadee is the Schema for the aadee API
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Aadee struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AadeeSpec   `json:"spec"`
	Status AadeeStatus `json:"status,omitempty"`
}

// AadeeList contains a list of Aadee
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AadeeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Aadee `json:"items"`
}
