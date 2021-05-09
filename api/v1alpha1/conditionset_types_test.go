package v1alpha1

import (
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"testing"
)

func TestConditionSet_Validate(t *testing.T) {
	const TaintEffectNotSupported v1.TaintEffect = "NotSupported"

	tests := []struct {
		Scenario     string
		ConditionSet ConditionSet
		Expected     bool
		ExpectedErr  bool
	}{
		{
			Scenario: "valid conditionset",
			ConditionSet: ConditionSet{
				Spec: ConditionSetSpec{
					Type:     "KubeletContainerRuntimeUnhealthy",
					Effect:   v1.TaintEffectNoExecute,
					TaintKey: "node.stakater.com/KubeletContainerRuntimeUnhealthy",
					Conditions: []NodeCondition{
						{
							Type:   KubeletUnhealthy,
							Status: v1.ConditionTrue,
						},
						{
							Type:   ContainerRuntimeUnhealthy,
							Status: v1.ConditionUnknown,
						},
					},
				},
			},
			Expected:    true,
			ExpectedErr: false,
		},
		{
			Scenario: "invalid conditionset",
			ConditionSet: ConditionSet{
				Spec: ConditionSetSpec{
					Type:     "KernelDeadlock",
					Effect:   v1.TaintEffectNoExecute,
					TaintKey: "node.stakater.com/KernelDeadlock",
					Conditions: []NodeCondition{
						{
							Type:   KernelDeadlock,
							Status: v1.ConditionFalse,
						},
					},
				},
			},
			Expected:    false,
			ExpectedErr: true,
		},
		{
			Scenario: "invalid taint effect",
			ConditionSet: ConditionSet{
				Spec: ConditionSetSpec{
					Type:     "KernelDeadlock",
					Effect:   TaintEffectNotSupported,
					TaintKey: "node.stakater.com/KernelDeadlock",
					Conditions: []NodeCondition{
						{
							Type:   KernelDeadlock,
							Status: v1.ConditionFalse,
						},
					},
				},
			},
			Expected:    false,
			ExpectedErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Scenario, func(t *testing.T) {
			status, err := tc.ConditionSet.Validate()
			if tc.Expected != status {
				t.Errorf("expected a valid ConditionSet")
			}

			if tc.ExpectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
