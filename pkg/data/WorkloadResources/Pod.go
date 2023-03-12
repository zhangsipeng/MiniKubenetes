package WorkloadResources

import (
	"example/Minik8s/pkg/data/ConfigAndStorageResources"
	"example/Minik8s/pkg/data/ObjectMeta"
	"time"
)

type Pod struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string
	Metadata   ObjectMeta.ObjectMeta
	Spec       PodSpec
	Status     PodStatus `yaml:"-"`
}

type PodSpec struct {
	IP            string `yaml:"-"`
	Container     []Container
	Volumes       []ConfigAndStorageResources.Volume
	NodeName      string `yaml:"nodeName"`
	RestartPolicy string `yaml:"restartPolicy"`
}

type PodStatus struct {
	StartTime time.Time
	Phase     string
	HostIP    string
	PodIP     string
	DockerIP  string
}

type ContainerPort struct {
	Protocal      string
	ContainerPort uint16 `yaml:"containerPort"`
	HostPort      uint16 `yaml:"hostPort"`
}

type Container struct {
	Name         string
	Image        string
	Command      []string
	Ports        []ContainerPort
	VolumeMounts []VolumeMount `yaml:"volumeMounts"`
}

type PodTemplateSpec struct {
	Metadata ObjectMeta.ObjectMeta
	Spec     PodSpec
}

type VolumeMount struct {
	MountPath string `yaml:"mountPath"`
	Name      string
}
