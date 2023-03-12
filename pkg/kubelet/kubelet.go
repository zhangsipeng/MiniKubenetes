package kubelet

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"example/Minik8s/pkg/apiclient"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"example/Minik8s/pkg/kubelet/dockerOp"
	"example/Minik8s/pkg/kubelet/joinAPIServer"
	"example/Minik8s/pkg/kubelet/reportStatus/reportNodeStatus"
	"example/Minik8s/pkg/kubelet/watchAPIServer"
	watchdockerevent "example/Minik8s/pkg/kubelet/watchDockerEvent"
	"fmt"
	"log"

	"github.com/docker/docker/client"
)

func StartService() {
	info := apiclient.GetInitInfo()
	runtimeConfig := apiclient.InitRuntimeConfig(info, setOthers)

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	nodeId := runtimeConfig.YamlConfig.Others["nodeId"]
	nodeByte := apiclient.Request(runtimeConfig,
		fmt.Sprintf("/api/v1/nodes/%s", nodeId), nil, "GET")
	if !info.Init {
		if string(nodeByte) == "null" {
			panic(errors.New("node does not exist in API Server"))
		}
	} else {
		if string(nodeByte) != "null" {
			panic(errors.New("node already exists in API Server"))
		}
		log.Printf("node %s: init\n", nodeId)
		joinAPIServer.JoinAPIServer(runtimeConfig)
		log.Printf("node %s: join cluster\n", nodeId)
		node, err := joinAPIServer.WaitForScheduler(runtimeConfig)
		if err != nil {
			panic(err)
		}
		log.Printf("node %s: get scheduled\n", nodeId)

		if err := dockerOp.ContainerRemoveAllK8s(cli); err != nil {
			panic(err)
		}
		if err := dockerOp.RemoveDockerNetworkIfExist(cli); err != nil {
			panic(err)
		}
		if err := dockerOp.CreateDockerNetwork(cli, *node); err != nil {
			panic(err)
		}
	}
	go reportNodeStatus.HeartBeat(runtimeConfig)
	go watchAPIServer.WatchPods(cli, runtimeConfig)
	go watchAPIServer.WatchJobs(runtimeConfig)
	watchdockerevent.WatchDockerEvent(cli, runtimeConfig)
}

func setOthers(yamlConfig *runtimedata.YamlConfig) {
	rand_byte := make([]byte, joinAPIServer.NODE_NAME_BYTE)
	_, err := rand.Read(rand_byte)
	if err != nil {
		panic(err)
	}
	yamlConfig.Others["nodeId"] = base64.RawURLEncoding.EncodeToString(rand_byte)
}
