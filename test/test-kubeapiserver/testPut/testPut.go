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
	pod := WorkloadResources.Pod{
		ApiVersion: "v1",
		Kind:       "pod",
		Metadata: ObjectMeta.ObjectMeta{
			Name:   "nginx1",
			Labels: map[string]string{"test": "1"},
		},
		Spec: WorkloadResources.PodSpec{
			Container: []WorkloadResources.Container{
				{
					Name:  "nginx",
					Image: "nginx:1.14.2",
					Ports: []WorkloadResources.ContainerPort{
						{
							Protocal:      "tcp",
							ContainerPort: 80,
						},
					},
				},
			},
		},
	}

	response := apiclient.Request(runtimeConfig, "/api/v1/pods", pod, "POST")

	log.Println(string(response))
}
