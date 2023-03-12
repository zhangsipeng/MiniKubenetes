package main

import (
	"encoding/json"
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/data/WorkloadResources"
)

func main() {
	info := apiclient.GetInitInfo()
	runtimeConfig, _ := apiclient.InitRuntimeConfigOrError(info, nil)
	relateUrl := "/api/v1/pods/"
	resBody := apiclient.Request(runtimeConfig, relateUrl, nil, "GET")
	podList := make([]WorkloadResources.Pod, 0)
	json.Unmarshal(resBody, &podList)
	for _, pod := range podList {
		println(pod.Metadata.Name)
	}

	relateUrl = "/api/v1/pods/nginx"
	resBody = apiclient.Request(runtimeConfig, relateUrl, nil, "GET")
	var pod WorkloadResources.Pod
	json.Unmarshal(resBody, &pod)
	println(pod.Metadata.Name)
}
