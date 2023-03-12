package main

import (
	"encoding/json"
	"example/Minik8s/pkg/data/ObjectMeta"
	"example/Minik8s/pkg/data/WorkloadResources"
)

func main() {
	pod := WorkloadResources.Pod{
		ApiVersion: "v1",
		Kind:       "pod",
		Metadata: ObjectMeta.ObjectMeta{
			Name: "nginx1",
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
	podJSON, _ := json.Marshal(pod)

	name, err := getName(podJSON)
	if err != nil {
		println(err)
		return
	}
	println(name)
}

func getName(body []byte) (string, error) {
	var bodyJson interface{}
	err := json.Unmarshal(body, &bodyJson)
	if err != nil {
		return "", err
	}

	name := ""
	for k, v := range bodyJson.(map[string]interface{}) {
		if k == "Metadata" {
			for k1, v1 := range v.(map[string]interface{}) {
				if k1 == "Name" {
					name = v1.(string)
					break
				}
			}
			break
		}
	}

	return name, nil
}
