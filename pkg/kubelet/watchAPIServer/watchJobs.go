package watchAPIServer

import (
	"encoding/json"
	"example/Minik8s/pkg/apiclient"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"example/Minik8s/pkg/data/WorkloadResources"
	"example/Minik8s/pkg/gpu"
	"example/Minik8s/pkg/kubeapiserver/watch"
	"fmt"
	"log"
	"strings"
)

func operateJob(event watch.WatchEvent, runtimeConfig runtimedata.RuntimeConfig) {
	eventKeyParts := strings.Split(event.Key, "/")
	jobName := eventKeyParts[len(eventKeyParts)-1]
	log.Printf("job %s: event %s", jobName, event.Type)
	switch event.Type {
	case "PUT":
		var job WorkloadResources.GPUJob
		if err := json.Unmarshal(
			apiclient.Request(runtimeConfig,
				fmt.Sprintf("/api/v1/gpujobs/%s", jobName),
				nil, "GET"),
			&job); err != nil {
			panic(err)
		}
		go gpu.SubmitJob(job, runtimeConfig)

		break
	}
}

func WatchJobs(runtimeConfig runtimedata.RuntimeConfig) {
	event := make(chan watch.WatchEvent)
	targetURL := fmt.Sprintf("/api/v1/jobInNode/%s/jobs/",
		runtimeConfig.YamlConfig.Others["nodeId"])
	log.Println(targetURL)
	go apiclient.WatchAPIWithRelativePath(runtimeConfig, targetURL, event)
	for {
		operateJob(<-event, runtimeConfig)
	}
}
