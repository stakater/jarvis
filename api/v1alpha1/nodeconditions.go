package v1alpha1

import v1 "k8s.io/api/core/v1"

// Supported Node conditions
const (

	// NodeReady means kubelet is healthy and ready to accept pods.
	NodeReady v1.NodeConditionType = "Ready"

	// NodeMemoryPressure means the kubelet is under pressure due to insufficient available memory.
	NodeMemoryPressure v1.NodeConditionType = "MemoryPressure"

	// NodeDiskPressure means the kubelet is under pressure due to insufficient available disk.
	NodeDiskPressure v1.NodeConditionType = "DiskPressure"

	// NodePIDPressure means the kubelet is under pressure due to insufficient available PID.
	NodePIDPressure v1.NodeConditionType = "PIDPressure"

	// NodeNetworkUnavailable means that network for the node is not correctly configured.
	NodeNetworkUnavailable v1.NodeConditionType = "NetworkUnavailable"

	// NTPProblem means that NTP Service is down.
	NTPProblem v1.NodeConditionType = "NTPProblem"

	// CorruptDockerOverlay2 means docker overlay2 is not functioning properly
	CorruptDockerOverlay2 v1.NodeConditionType = "CorruptDockerOverlay2"

	// ContainerRuntimeUnhealthy means Container runtime on the node is not functioning properly
	ContainerRuntimeUnhealthy v1.NodeConditionType = "ContainerRuntimeUnhealthy"

	// KubeletUnhealthy means kubelet on the node is not functioning properly
	KubeletUnhealthy v1.NodeConditionType = "KubeletUnhealthy"

	// KernelDeadlock means kernel might have deadlock
	KernelDeadlock v1.NodeConditionType = "KernelDeadlock"

	// ReadonlyFilesystem means Filesystem is read-only"
	ReadonlyFilesystem v1.NodeConditionType = "ReadonlyFilesystem"

	// FrequentUnregisterNetDevice means node is not functioning properly
	FrequentUnregisterNetDevice v1.NodeConditionType = "FrequentUnregisterNetDevice"

	// FrequentKubeletRestart means kubelet is not functioning properly
	FrequentKubeletRestart v1.NodeConditionType = "FrequentKubeletRestart"

	// FrequentDockerRestart means docker is not functioning properly
	FrequentDockerRestart v1.NodeConditionType = "FrequentDockerRestart"

	// FrequentContainerdRestart means containerd is not functioning properly
	FrequentContainerdRestart v1.NodeConditionType = "FrequentContainerdRestart"
)

// ValidNodeConditionStatusMapping provides valid supported condition status against
// the specific node condition type that could be used in creation of NodeConditionSet.
// This will be used in validating list of NodeConditionSet's Conditions.
// This is a precautionary measure to make sure that we don't allow any NodeConditions
// that could lead in healing a already healthy node.
// For example:
// Suppose one of the Condition belonging to a ConditionSet is :
//  NodeConditionType: NodeMemoryPressure and Status: False
// Where status False meaning nodes which don't have memory pressure, which all healthy node
// will have.
// This could lead to an unnecessary healing enforcement on a healthy nodes.
// To avoid such case, we will validate NodeConditionSet's NodeConditions at the time of creation.
var ValidNodeConditionStatusMapping = map[v1.NodeConditionType][]v1.ConditionStatus{
	NodeReady:                   []v1.ConditionStatus{v1.ConditionFalse, v1.ConditionUnknown},
	NodeMemoryPressure:          []v1.ConditionStatus{v1.ConditionTrue, v1.ConditionUnknown},
	NodeDiskPressure:            []v1.ConditionStatus{v1.ConditionTrue, v1.ConditionUnknown},
	NodePIDPressure:             []v1.ConditionStatus{v1.ConditionTrue, v1.ConditionUnknown},
	NodeNetworkUnavailable:      []v1.ConditionStatus{v1.ConditionTrue, v1.ConditionUnknown},
	NTPProblem:                  []v1.ConditionStatus{v1.ConditionTrue, v1.ConditionUnknown},
	CorruptDockerOverlay2:       []v1.ConditionStatus{v1.ConditionTrue, v1.ConditionUnknown},
	ContainerRuntimeUnhealthy:   []v1.ConditionStatus{v1.ConditionTrue, v1.ConditionUnknown},
	KubeletUnhealthy:            []v1.ConditionStatus{v1.ConditionTrue, v1.ConditionUnknown},
	KernelDeadlock:              []v1.ConditionStatus{v1.ConditionTrue, v1.ConditionUnknown},
	ReadonlyFilesystem:          []v1.ConditionStatus{v1.ConditionTrue, v1.ConditionUnknown},
	FrequentUnregisterNetDevice: []v1.ConditionStatus{v1.ConditionTrue, v1.ConditionUnknown},
	FrequentKubeletRestart:      []v1.ConditionStatus{v1.ConditionTrue, v1.ConditionUnknown},
	FrequentDockerRestart:       []v1.ConditionStatus{v1.ConditionTrue, v1.ConditionUnknown},
	FrequentContainerdRestart:   []v1.ConditionStatus{v1.ConditionTrue, v1.ConditionUnknown},
}
