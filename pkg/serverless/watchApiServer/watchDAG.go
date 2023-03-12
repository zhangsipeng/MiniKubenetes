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

func addServerlessDAG(event watch.WatchEvent) {
	switch event.Type {
	case "PUT":
		var serverlessDAG Serverless.DAG
		err := json.Unmarshal(([]byte)(event.Value), &serverlessDAG)
		if err != nil {
			log.Println(err)
			return
		}
		if serverlessDAG.Status.Phase != "init" {
			return
		}
		log.Println("notice new DAG " + serverlessDAG.Metadata.Name)
		go server.StartServerlessDAG(serverlessDAG)
	}
}

func WatchServerlessDAG() {
	event := make(chan watch.WatchEvent)
	go apiclient.WatchAPIWithRelativePath(data.GetCredential(), "/api/v1/serverlessDAG/", event)
	for {
		addServerlessDAG(<-event)
	}
}
