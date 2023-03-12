package main

import (
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/data/ObjectMeta"
	"example/Minik8s/pkg/data/ServiceResources"
)

func main() {
	info := apiclient.GetInitInfo()
	config := apiclient.InitRuntimeConfig(info, nil)
	service := ServiceResources.Service{
		ApiVersion: "v1",
		Kind:       "service",
		Metadata: ObjectMeta.ObjectMeta{
			Name: "test-service",
		},
		Spec: ServiceResources.ServiceSpec{
			ClusterIP: "192.168.100.101",
			Ports: []ServiceResources.ServicePort{
				{Port: 10000, TargetPort: 80},
			},
			Selector: map[string]string{"test": "1"},
		},
		Status: ServiceResources.ServiceStatus{},
	}
	apiclient.Request(config, "/api/v1/service", service, "POST")
}
