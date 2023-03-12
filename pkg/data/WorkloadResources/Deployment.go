package WorkloadResources

import "example/Minik8s/pkg/data/ObjectMeta"

type Deployment struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string
	Metadata   ObjectMeta.ObjectMeta
	Spec       DeploymentSpec
	Status     DeploymentStatus `yaml:"-"`
}

type DeploymentSpec struct {
	Replicas int32
	Selector map[string]string
	Template PodTemplateSpec
}

type DeploymentStatus struct {
	Replicas          int32
	AvailableReplicas int32
	ReadyReplicas     int32
	PodName           []string
}
