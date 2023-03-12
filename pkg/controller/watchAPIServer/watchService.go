package watchAPIServer

import (
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/controller/AddEndpoints"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"example/Minik8s/pkg/kubeapiserver/watch"
	"fmt"
)

func WatchService(credential runtimedata.RuntimeConfig) {
	event := make(chan watch.WatchEvent)
	targetURL := fmt.Sprintf("/api/v1/service")
	go apiclient.WatchAPIWithRelativePath(credential, targetURL, event)
	for {
		AddEndpoints.AddEndpoints(credential, <-event)
	}
}
