package main

import (
	"example/Minik8s/pkg/apiclient"
)

func main() {
	info := apiclient.GetInitInfo()
	runtimeConfig, _ := apiclient.InitRuntimeConfigOrError(info, nil)
	relateUrl := "/api/v1/pods/nginx1"
	apiclient.Request(runtimeConfig, relateUrl, nil, "DELETE")
}
