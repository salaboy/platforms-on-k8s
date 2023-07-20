package v1alpha1

//go:generate controller-gen object paths=$GOFILE

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type Parameters struct {
	InstallInfra bool `json:"installInfra"`
}

type CompositionSelector struct {
	MatchLabels map[string]string `json:"matchLabels"`
}

type EnvironmentSpec struct {
	Parameters          Parameters          `json:"parameters"`
	CompositionSelector CompositionSelector `json:"compositionSelector,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Environment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec EnvironmentSpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type EnvironmentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Environment `json:"items"`
}
