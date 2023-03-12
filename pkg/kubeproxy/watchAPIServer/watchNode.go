package watchAPIServer

import (
	"encoding/json"
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/data/ClusterResources"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"example/Minik8s/pkg/kubeapiserver/watch"
	"example/Minik8s/pkg/kubeproxy/vxlan"
	"log"
	"strings"
)

func handleNodeEvent(event watch.WatchEvent) {
	// ignore if it is about things like "nodes/NODE/pods/POD"
	if strings.Index(event.Key[len("nodes/"):], "/") != -1 {
		return
	}
	switch event.Type {
	case "PUT": // new node
		var node ClusterResources.Node
		if err := json.Unmarshal([]byte(event.Value), &node); err != nil {
			log.Printf("warning: %s", err.Error())
			return
		}
		if node.Status.Phase == "init" { // no CIDR allocated
			return
		}
		if err := vxlan.AddPeer(node.Status.Addresses[0].Address,
			node.Spec.NodeVxlanCIDR, node.Spec.PodCIDR); err != nil {
			log.Printf("warning: %s", err.Error())
			return
		}
		break
	case "DELETE":
		// TODO
	}
}

func WatchNode(runtimeConfig runtimedata.RuntimeConfig) {
	event := make(chan watch.WatchEvent)
	relativeUrlPath := "/api/v1/nodes/"
	go apiclient.WatchAPIWithRelativePath(runtimeConfig, relativeUrlPath, event)
	for {
		handleNodeEvent(<-event)
	}
}
