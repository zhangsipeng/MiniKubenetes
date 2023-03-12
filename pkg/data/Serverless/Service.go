package Serverless

import "example/Minik8s/pkg/data/ObjectMeta"

type Service struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string
	Metadata   ObjectMeta.ObjectMeta
	Spec       ServiceSpec
	Status     ServiceStatus
}

type ServiceSpec struct {
	GitUrl      string `yaml:"gitUrl"`
	MaxReplicas int32  `yaml:"maxReplicas"`
	Input       []string
	Output      []string
}

type ServiceStatus struct {
	Phase    string
	Replicas int32
}
