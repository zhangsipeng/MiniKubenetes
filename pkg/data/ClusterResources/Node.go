package ClusterResources

import (
	"example/Minik8s/pkg/data/ObjectMeta"
	"time"
)

type Node struct {
	ApiVersion string
	Kind       string
	Metadata   ObjectMeta.ObjectMeta
	Spec       NodeSpec
	Status     NodeStatus
}

type NodeSpec struct {
	NodeVxlanCIDR string
	PodCIDR       string
	PodCIDRs      []string
	Unschedulable bool
}

type NodeStatus struct {
	Addresses       []NodeAddress
	Conditions      []NodeCondition
	Image           []ContainerImage
	NodeInfo        NodeSystemInfo
	Phase           string
	VolumesAttached []AttachedVolume
	VolumesInUse    []string
}

type NodeAddress struct {
	Address string
	Type    string
}

type NodeCondition struct {
	Status             string
	Type               string
	LastHeartBeatTime  time.Time
	LastTransitionTime time.Time
}

type ContainerImage struct {
	Name      []string
	SizeBytes int64
}

type NodeSystemInfo struct {
	Architecture            string
	BootId                  string
	ContainerRuntimeVersion string
	KernelVersion           string
	KubeProxyVersion        string
	KubeletVersion          string
	MachineID               string
	OperatingSystem         string
	OsImage                 string
	SystemUUID              string
}

type AttachedVolume struct {
	DevicePath string
	Name       string
}
