package watchAPIServer

import (
	"encoding/json"
	"example/Minik8s/pkg/apiclient"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"example/Minik8s/pkg/data/ServiceResources"
	"example/Minik8s/pkg/kubeapiserver/watch"
	"example/Minik8s/pkg/kubeproxy/iptables"
	"log"
)

func addEndpoints(ipt *iptables.IPTable, service ServiceResources.Service, endpoints ServiceResources.Endpoints) (err error) {
	if len(endpoints.PodSet) > 0 {
		log.Printf("add endpoint %s of service %s\n",
			endpoints.Metadata.Name,
			service.Metadata.Name)
		err = ipt.AddService(service, endpoints.PodSet)
	}
	return
}

func changeEndpoints(ipt *iptables.IPTable, service ServiceResources.Service, event watch.WatchEvent) {
	switch event.Type {
	case "PUT":
		var endpoints ServiceResources.Endpoints
		err := json.Unmarshal(([]byte)(event.Value), &endpoints)
		if err != nil {
			panic(err)
		}
		err = addEndpoints(ipt, service, endpoints)
		if err != nil {
			panic(err)
		}
	}
}

func watchEndpoints(ipt *iptables.IPTable, service ServiceResources.Service, credential runtimedata.RuntimeConfig) {
	event := make(chan watch.WatchEvent)
	relativeUrlPath := "/api/v1/endpoints/" + service.Metadata.Name + "-endpoints/"
	{
		var endpoints ServiceResources.Endpoints
		err := json.Unmarshal(
			apiclient.Request(credential, relativeUrlPath, nil, "GET"),
			&endpoints)
		if err != nil {
			panic(err)
		}
		addEndpoints(ipt, service, endpoints)
	}
	go apiclient.WatchAPIWithRelativePath(credential, relativeUrlPath, event)
	for {
		changeEndpoints(ipt, service, <-event)
	}
}

func addService(ipt *iptables.IPTable, credential runtimedata.RuntimeConfig, event watch.WatchEvent) {
	switch event.Type {
	case "PUT":
		var service ServiceResources.Service
		err := json.Unmarshal(([]byte)(event.Value), &service)
		if err != nil {
			log.Println(err)
			return
		}
		go watchEndpoints(ipt, service, credential)
	}
}

func WatchService(ipt *iptables.IPTable, credential runtimedata.RuntimeConfig) {
	event := make(chan watch.WatchEvent)
	relativeUrlPath := "/api/v1/service/"
	go apiclient.WatchAPIWithRelativePath(credential, relativeUrlPath, event)
	for {
		addService(ipt, credential, <-event)
	}
}
