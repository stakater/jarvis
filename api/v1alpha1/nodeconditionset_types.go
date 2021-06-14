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
	"github.com/pkg/errors"
	"github.com/stakater/jarvis/utils/slice"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type NodeCondition struct {
	// Type of node condition
	// +required
	// +kubebuilder:validation:Required
	Type v1.NodeConditionType `json:"type"`

	// Status of the condition, one of True, False, Unknown
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=True;False;Unknown
	Status v1.ConditionStatus `json:"status"`
}

type NodeConditionSetName string

// NodeConditionSetSpec defines the desired state of NodeConditionSet
type NodeConditionSetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Type of NodeConditionSet
	// +required
	// +kubebuilder:validation:Required
	Name NodeConditionSetName `json:"name"`

	// Effect indicates the taint effect to match.
	// Valid effects are NoSchedule, PreferNoSchedule, and NoExecute
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=NoSchedule;PreferNoSchedule;NoExecute
	Effect v1.TaintEffect `json:"effect"`

	// The taint key to be applied to a node
	// +required
	// +kubebuilder:validation:Required
	TaintKey string `json:"taintKey"`

	// NodeConditions is an array of unique NodeCondition, that collectively define a NodeConditionSet
	// +required
	// +kubebuilder:validation:Required
	NodeConditions []NodeCondition `json:"nodeConditions"`
}

// NodeConditionSetStatus defines the observed state of NodeConditionSet
type NodeConditionSetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Error
	// +kubebuilder:validation:Optional
	Error string `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:categories=all;auto-healer,shortName=ncs

// NodeConditionSet is the Schema for the NodeConditionSets API
type NodeConditionSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeConditionSetSpec   `json:"spec,omitempty"`
	Status NodeConditionSetStatus `json:"status,omitempty"`
}

func (ncs *NodeConditionSet) Validate() (bool, error) {

	switch ncs.Spec.Effect {
	case v1.TaintEffectPreferNoSchedule, v1.TaintEffectNoSchedule, v1.TaintEffectNoExecute:
	default:
		return false, errors.New("invalid Taint effect")
	}

	for _, nodeCondition := range ncs.Spec.NodeConditions {
		validStatus, ok := ValidNodeConditionStatusMapping[nodeCondition.Type]
		if ok {
			if !slice.Contains(validStatus, nodeCondition.Status) {
				return false, errors.New("invalid NodeCondition status")
			}
		} else {
			return false, errors.New("unsupported NodeConditionType")
		}
	}

	return true, nil
}

//+kubebuilder:object:root=true

// NodeConditionSetList contains a list of NodeConditionSet
type NodeConditionSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeConditionSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodeConditionSet{}, &NodeConditionSetList{})
}
