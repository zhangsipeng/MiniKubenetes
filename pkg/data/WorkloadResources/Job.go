package WorkloadResources

import (
	"example/Minik8s/pkg/data/ObjectMeta"
)

type GPUJob struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string
	Metadata   ObjectMeta.ObjectMeta
	Giturl     string
	Config     GPUConfig
	Res        string
	Phase      string
	NodeName   string `yaml:"-"`
}

type GPUConfig struct {
	Nodes      int
	Tskpernode int
	Cpupertsk  int
	Gpunum     int
}
