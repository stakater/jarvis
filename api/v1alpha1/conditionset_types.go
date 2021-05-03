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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type NodeConditionType string

type ConditionStatus string

type NodeCondition struct {
	// Type of node condition
	// +required
	// +kubebuilder:validation:Required
	Type NodeConditionType `json:"type"`

	// Status of the condition, one of True, False, Unknown
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=True;False;Unknown
	Status ConditionStatus `json:"status"`
}

type ConditionSetType string

type TaintEffect string

// ConditionSetSpec defines the desired state of ConditionSet
type ConditionSetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Type of ConditionSet
	// +required
	// +kubebuilder:validation:Required
	Type ConditionSetType `json:"type"`

	// Effect indicates the taint effect to match.
	// Valid effects are NoSchedule, PreferNoSchedule, and NoExecute
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=NoSchedule;PreferNoSchedule;NoExecute
	Effect TaintEffect `json:"effect"`

	// The taint key to be applied to a node
	// +required
	// +kubebuilder:validation:Required
	TaintKey string `json:"taintKey"`

	// Conditions is an array of unique NodeConditions, that collectively define a ConditionSet
	// +required
	// +kubebuilder:validation:Required
	Conditions []NodeCondition `json:"conditions"`
}

// ConditionSetStatus defines the observed state of ConditionSet
type ConditionSetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Error
	// +kubebuilder:validation:Optional
	Error string `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:categories=all;auto-healer,shortName=cs

// ConditionSet is the Schema for the conditionsets API
type ConditionSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConditionSetSpec   `json:"spec,omitempty"`
	Status ConditionSetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ConditionSetList contains a list of ConditionSet
type ConditionSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConditionSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConditionSet{}, &ConditionSetList{})
}
