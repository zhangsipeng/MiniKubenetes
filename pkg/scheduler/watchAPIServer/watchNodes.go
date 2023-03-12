package watchAPIServer

import (
	"encoding/json"
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/data/ClusterResources"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"example/Minik8s/pkg/kubeapiserver/watch"
	"example/Minik8s/pkg/scheduler/data"
	"log"
)

func addNode(runtimeConfig runtimedata.RuntimeConfig, event watch.WatchEvent) {
	switch event.Type {
	case "PUT":
		var node ClusterResources.Node
		err := json.Unmarshal(([]byte)(event.Value), &node)
		if err != nil {
			log.Println(err)
			return
		}
		if node.Kind != "node" {
			return
		}
		if node.Status.Phase != "init" {
			return
		}
		log.Println(event.Key)
		data.AddNode(&node)
		apiclient.Request(runtimeConfig, "/api/v1/nodes/", node, "PUT")
	}
}

func WatchNodes(credential runtimedata.RuntimeConfig) {
	event := make(chan watch.WatchEvent)
	go apiclient.WatchAPIWithRelativePath(credential, "/api/v1/nodes/", event)
	for {
		addNode(credential, <-event)
	}
}
