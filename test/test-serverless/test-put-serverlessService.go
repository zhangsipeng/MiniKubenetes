package main

import (
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/data/ObjectMeta"
	"example/Minik8s/pkg/data/Serverless"
	"log"
)

func main() {
	serverlessBuild := Serverless.Service{
		ApiVersion: "v1",
		Kind:       "serverlessService",
		Metadata: ObjectMeta.ObjectMeta{
			Name: "add",
		},
		Spec: Serverless.ServiceSpec{
			GitUrl:      "/root/serverless/add",
			MaxReplicas: 5,
			Input:       []string{"A", "B"},
			Output:      []string{"A"},
		},
		Status: Serverless.ServiceStatus{
			Phase: "init",
		},
	}

	info := apiclient.GetInitInfo()
	runtimeConfig, _ := apiclient.InitRuntimeConfigOrError(info, nil)

	response := apiclient.Request(runtimeConfig, "/api/v1/serverless/", serverlessBuild, "POST")

	log.Println(string(response))
}
