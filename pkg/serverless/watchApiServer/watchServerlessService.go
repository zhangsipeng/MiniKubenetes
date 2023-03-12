package watchApiServer

import (
	"encoding/json"
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/data/Serverless"
	"example/Minik8s/pkg/kubeapiserver/watch"
	"example/Minik8s/pkg/serverless/data"
	"example/Minik8s/pkg/serverless/server"
	"log"
)

func addServerlessService(event watch.WatchEvent) {
	switch event.Type {
	case "PUT":
		var serverlessService Serverless.Service
		err := json.Unmarshal(([]byte)(event.Value), &serverlessService)
		if err != nil {
			log.Println(err)
			return
		}
		if serverlessService.Status.Phase == "removing" {
			server.DeleteServerlessService(serverlessService, false)
			return
		}
		if serverlessService.Status.Phase == "change" {
			server.UpdateServerlessService(serverlessService)
			return
		}
		if serverlessService.Status.Phase != "init" {
			return
		}
		log.Println("notice new Service " + serverlessService.Metadata.Name)
		go server.StartServerlessService(serverlessService)
	}
}

func WatchServerlessService() {
	event := make(chan watch.WatchEvent)
	go apiclient.WatchAPIWithRelativePath(data.GetCredential(), "/api/v1/serverless/", event)
	for {
		addServerlessService(<-event)
	}
}
