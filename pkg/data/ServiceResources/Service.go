package ServiceResources

import (
	"example/Minik8s/pkg/data/ObjectMeta"
	"time"
)

type Service struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string
	Metadata   ObjectMeta.ObjectMeta
	Spec       ServiceSpec
	Status     ServiceStatus `yaml:"-"`
}

type ServiceSpec struct {
	ClusterIP string `yaml:"clusterIP"`
	Selector  map[string]string
	Ports     []ServicePort
	Type      string
}

type ServiceStatus struct {
	Conditions   []Condition
	LoadBalancer LoadBalancerStatus
}

type ServicePort struct {
	Port       int32
	TargetPort int32 `yaml:"targetPort"`
	Protocol   string
	Name       string
	NodePort   int32 `yaml:"nodePort"`
}

type Condition struct {
	LastTransitionTime time.Time
	Message            string
	Reason             string
	Status             string
	Type               string
	ObservedGeneration int64
}

type LoadBalancerStatus struct {
	Ingress []LoadBalancerIngress
}

type LoadBalancerIngress struct {
	Hostname string
	Ip       string
	Ports    []PortStatus
}

type PortStatus struct {
	Port     int32
	Protocol string
	Error    string
}
