package scheduler

import (
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/scheduler/watchAPIServer"
)

func StartService() {
	info := apiclient.GetInitInfo()
	runtimeConfig := apiclient.InitRuntimeConfig(info, nil)
	go watchAPIServer.WatchNodes(runtimeConfig)
	go watchAPIServer.WatchPods(runtimeConfig)
	watchAPIServer.WatchJobs(runtimeConfig)
}
