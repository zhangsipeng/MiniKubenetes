package schedule

import (
	"encoding/json"
	"example/Minik8s/pkg/apiclient"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"example/Minik8s/pkg/data/WorkloadResources"
	"example/Minik8s/pkg/kubeapiserver/watch"
	"fmt"
	"log"
)

func ScheduleJobs(credential runtimedata.RuntimeConfig, event watch.WatchEvent) {
	switch event.Type {
	case "PUT":
		var job WorkloadResources.GPUJob
		err := json.Unmarshal(([]byte)(event.Value), &job)
		if err != nil {
			log.Println(err)
			return
		}
		if job.Phase != "init" {
			return
		}
		node := selectNode()
		if node == nil {
			log.Println("no available node!")
			return
		}
		if err != nil {
			log.Println(err)
			return
		}

		job.Phase = "scheduled"
		job.NodeName = node.Metadata.Name
		apiclient.Request(credential,
			"/api/v1/gpujobs/", job, "PUT")
		apiclient.Request(credential,
			fmt.Sprintf("/api/v1/jobInNode/%s/jobs/", node.Metadata.Name),
			job, "POST")
	}
}
