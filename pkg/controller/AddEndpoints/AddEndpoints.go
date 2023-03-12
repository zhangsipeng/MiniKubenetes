package AddEndpoints

import (
	"encoding/json"
	"example/Minik8s/pkg/apiclient"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"example/Minik8s/pkg/data/ServiceResources"
	"example/Minik8s/pkg/data/WorkloadResources"
	"example/Minik8s/pkg/kubeapiserver/watch"
	"fmt"
)

func JudgeLabel(podlabel map[string]string, servicelabel map[string]string) bool {
	for k, v := range podlabel {
		v_service, ok := servicelabel[k]
		if ok && v_service == v {
			continue
		} else {
			return false
		}
	}
	return true
}
func AddEndpoints(credential runtimedata.RuntimeConfig, event watch.WatchEvent) {
	switch event.Type {
	case "PUT":
		{
			var service ServiceResources.Service
			err := json.Unmarshal(([]byte)(event.Value), &service)
			if err != nil {
				panic("err decode service json")
			}
			var endpoints ServiceResources.Endpoints
			endpoints.Kind = "endpoints"
			endpoints.Metadata.Name = service.Metadata.Name + "-endpoints"
			response := apiclient.Request(credential, "/api/v1/pods", nil, "GET")
			str_response := string(response[:])
			if str_response == "null" {
				fmt.Printf("pod for service unavailable now,plz try to create pod and recreate service")
				return
			} else {
				podSelectNum := 0
				var podlist []WorkloadResources.Pod
				err := json.Unmarshal(response, &podlist)
				if err != nil {
					panic("Fatel:decode podlist json failed \n")
				}
				total := len(podlist)
				for i := 0; i < total; i++ {
					item := podlist[i]
					if JudgeLabel(item.Metadata.Labels, service.Spec.Selector) {
						podSelectNum++
						endpoints.PodSet = append(endpoints.PodSet, item)
					}
				}
				if podSelectNum == 0 {
					fmt.Printf("pod for service unavailable now,plz try to create pod and recreate service")
					return
				} else {
					response = apiclient.Request(credential, "/api/v1/endpoints", endpoints, "POST")
					fmt.Printf("request to successfully create endpoints http return%s", response)
				}
			}

		}
	default:
		{
			fmt.Printf("Unsupported event!")
		}
	}

}
