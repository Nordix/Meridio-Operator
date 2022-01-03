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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// TrenchSpec defines the desired state of Trench
type TrenchSpec struct {
	// Defines the IP family of the trench. It should be set according to what type of traffic is expected in the trench.
	// Valid values:
	// - dualstack (default)
	// - ipv4
	// - ipv6
	// +optional
	Family string `json:"family,omitempty"`

	AnotherParam string `json:"another,omitempty"`
}

// TrenchStatus defines the observed state of Trench
type TrenchStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:storageversion
//+kubebuilder:printcolumn:name="family",type=string,JSONPath=`.spec.family`

// Trench is the Schema for the trenches API
type Trench struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TrenchSpec   `json:"spec,omitempty"`
	Status TrenchStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TrenchList contains a list of Trench
type TrenchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Trench `json:"items"`
}

func (r *Trench) GroupResource() schema.GroupResource {
	return schema.GroupResource{
		Group:    r.GroupVersionKind().Group,
		Resource: r.GroupVersionKind().Kind,
	}
}

func (r *Trench) GroupKind() schema.GroupKind {
	return schema.GroupKind{
		Group: r.GroupVersionKind().Group,
		Kind:  r.GroupVersionKind().Kind,
	}
}

func init() {
	SchemeBuilder.Register(&Trench{}, &TrenchList{})
}
