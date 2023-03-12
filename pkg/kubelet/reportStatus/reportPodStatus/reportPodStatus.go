package reportPodStatus

import (
	"encoding/json"
	"example/Minik8s/pkg/apiclient"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"example/Minik8s/pkg/data/WorkloadResources"
	"fmt"
)

func ReportPodStatus(podName string, phase string, runtimeConfig runtimedata.RuntimeConfig) {
	var pod WorkloadResources.Pod
	targetUrl := fmt.Sprintf("/api/v1/pods/%s/", podName)
	json.Unmarshal(
		apiclient.Request(runtimeConfig, targetUrl, nil, "GET"),
		&pod)
	pod.Status.Phase = phase
	apiclient.Request(runtimeConfig, "/api/v1/pods/", pod, "PUT")
}
