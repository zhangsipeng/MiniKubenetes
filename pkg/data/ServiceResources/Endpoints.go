package ServiceResources

import (
	"example/Minik8s/pkg/data/ObjectMeta"
	"example/Minik8s/pkg/data/WorkloadResources"
)

type Endpoints struct {
	ApiVersion string
	Kind       string
	Metadata   ObjectMeta.ObjectMeta
	PodSet     []WorkloadResources.Pod
}
