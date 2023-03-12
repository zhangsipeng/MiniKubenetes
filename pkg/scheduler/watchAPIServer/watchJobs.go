package watchAPIServer

import (
	"example/Minik8s/pkg/apiclient"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"example/Minik8s/pkg/kubeapiserver/watch"
	"example/Minik8s/pkg/scheduler/schedule"
)

func WatchJobs(credential runtimedata.RuntimeConfig) {
	event := make(chan watch.WatchEvent)
	go apiclient.WatchAPIWithRelativePath(credential, "/api/v1/gpujobs/", event)
	for {
		schedule.ScheduleJobs(credential, <-event)
	}
}
