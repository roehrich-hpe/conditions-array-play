/*
Copyright 2023.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Condition types used in []metav1.Conditions.  These are nouns.

	// Indicates whether the Beak CR is available
	BirdConditionBeakResource = "BeakResource"
)

const (
	// Condition reasons used in []metav1.Conditions.  These are past-tense
	// verbs containing the reason for the last transition.

	// Indicates the resource was created
	BirdConditionResourceCreated = "ResourceCreated"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BirdSpec defines the desired state of Bird
type BirdSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Bird. Edit bird_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// BirdStatus defines the observed state of Bird
type BirdStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// NOTE: The following Condition build tags are from the example in:
	//   k8s.io/apimachinery/pkg/api/meta/v1/types.go

	// Conditions represents the observations of Bird's current state
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Bird is the Schema for the birds API
type Bird struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BirdSpec   `json:"spec,omitempty"`
	Status BirdStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BirdList contains a list of Bird
type BirdList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Bird `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Bird{}, &BirdList{})
}
