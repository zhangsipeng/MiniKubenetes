package Serverless

import "example/Minik8s/pkg/data/ObjectMeta"

type DAG struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string
	Metadata   ObjectMeta.ObjectMeta
	Spec       DAGSpec
	Status     DAGStatus `yaml:"-"`
}

type DAGSpec struct {
	Steps []Step
	Input []string
}

type Step struct {
	Name   string
	Type   string
	Task   StepTask
	Choice StepChoice
}

type StepTask struct {
	Function Service
}

type StepChoice struct {
	Key  string
	Jump map[string]string
}

type DAGStatus struct {
	Phase string
}
