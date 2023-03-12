package watch

import (
	"encoding/json"
	"example/Minik8s/pkg/apiclient"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"example/Minik8s/pkg/data/ServiceResources"
	"example/Minik8s/pkg/dnsManager/dnsOp"
	"example/Minik8s/pkg/kubeapiserver/watch"
	"log"
)

func handleService(event watch.WatchEvent) (err error) {
	switch event.Type {
	case "PUT":
		var service ServiceResources.Service
		err = json.Unmarshal([]byte(event.Value), &service)
		if err != nil {
			return
		}
		if service.Spec.ClusterIP != "" {
			log.Printf("adding service %s\n", service.Metadata.Name)
			err = dnsOp.AddService(service)
		}
		return
	case "DELETE":
		// TODO
	}
	return nil
}

func WatchService(runtimeConfig runtimedata.RuntimeConfig) {
	event := make(chan watch.WatchEvent)
	relativeUrlPath := "/api/v1/service/"
	go apiclient.WatchAPIWithRelativePath(runtimeConfig, relativeUrlPath, event)
	for {
		handleService(<-event)
	}
}
