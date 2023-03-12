package main

import (
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/data/ObjectMeta"
	"example/Minik8s/pkg/data/ServiceResources"
	"example/Minik8s/pkg/data/WorkloadResources"
)

func main() {
	info := apiclient.GetInitInfo()
	config := apiclient.InitRuntimeConfig(info, nil)
	endpoints := ServiceResources.Endpoints{
		ApiVersion: "v1",
		Kind:       "endpoints",
		Metadata:   ObjectMeta.ObjectMeta{Name: "test-service-endpoints"},
		PodSet: []WorkloadResources.Pod{
			{Status: WorkloadResources.PodStatus{PodIP: "192.168.100.100"}},
			{Status: WorkloadResources.PodStatus{PodIP: "192.168.100.110"}}},
	}
	apiclient.Request(config, "/api/v1/endpoints", endpoints, "POST")
}
