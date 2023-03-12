package main

import (
	"example/Minik8s/pkg/apiclient"
	"log"
)

func main() {
	info := apiclient.GetInitInfo()
	runtimeConfig, _ := apiclient.InitRuntimeConfigOrError(info, nil)
	response := apiclient.Request(runtimeConfig, "/api/v1/serverless/add/", nil, "DELETE")
	log.Println(string(response))
}
