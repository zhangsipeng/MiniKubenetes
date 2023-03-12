package main

import (
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/data/ObjectMeta"
	"example/Minik8s/pkg/data/WorkloadResources"
	"log"
)

func main() {
	info := apiclient.GetInitInfo()
	runtimeConfig := apiclient.InitRuntimeConfig(info, nil)
	job := WorkloadResources.GPUJob{
		ApiVersion: "v1",
		Kind:       "job",
		Metadata:   ObjectMeta.ObjectMeta{Name: "test"},
	}

	response := apiclient.Request(runtimeConfig, "/api/v1/gpujobs", job, "POST")

	log.Println(string(response))
}
