package watch

import (
	"encoding/json"
	"example/Minik8s/pkg/apiclient"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"example/Minik8s/pkg/data/WorkloadResources"
	"example/Minik8s/pkg/dnsManager/dnsOp"
	"example/Minik8s/pkg/kubeapiserver/watch"
	"log"
)

func handlePod(event watch.WatchEvent) (err error) {
	switch event.Type {
	case "PUT":
		var pod WorkloadResources.Pod
		err = json.Unmarshal([]byte(event.Value), &pod)
		if err != nil {
			return
		}
		if pod.Spec.IP != "" {
			log.Printf("adding pod %s\n", pod.Metadata.Name)
			err = dnsOp.AddPod(pod)
		}
		return
	case "DELETE":
		// TODO
	}
	return nil
}

func WatchPod(runtimeConfig runtimedata.RuntimeConfig) {
	event := make(chan watch.WatchEvent)
	relativeUrlPath := "/api/v1/pods/"
	go apiclient.WatchAPIWithRelativePath(runtimeConfig, relativeUrlPath, event)
	for {
		handlePod(<-event)
	}
}
