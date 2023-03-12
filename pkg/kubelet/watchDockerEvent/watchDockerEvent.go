package watchdockerevent

import (
	"context"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"example/Minik8s/pkg/kubelet/reportStatus/reportPodStatus"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

func dockerWatcher(cli *client.Client, eventChan chan events.Message) {
	c, _ := cli.Events(context.TODO(), types.EventsOptions{})
	for {
		eventChan <- <-c
	}
}

func WatchDockerEvent(cli *client.Client, runtimeConfig runtimedata.RuntimeConfig) {
	eventChan := make(chan events.Message)
	go dockerWatcher(cli, eventChan)
	defer close(eventChan)
	for {
		event := <-eventChan
		if event.Type == events.ContainerEventType {
			if event.Action == "die" {
				if event.Actor.Attributes["restartPolicy"] != "always" {
					podName := event.Actor.Attributes["pod"]
					reportPodStatus.ReportPodStatus(podName, "died", runtimeConfig)
				}
			}
		}
	}
}
