package schedule

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/data/ClusterResources"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"example/Minik8s/pkg/data/WorkloadResources"
	"example/Minik8s/pkg/kubeapiserver/watch"
	"example/Minik8s/pkg/scheduler/data"
	"fmt"
	"log"
	"math/big"
	"net"
)

var podIndexOfNode = make(map[string]int)

const maxPodIndexOfNode = 253

func SchedulePods(credential runtimedata.RuntimeConfig, event watch.WatchEvent) {
	switch event.Type {
	case "PUT":
		var pod WorkloadResources.Pod
		err := json.Unmarshal(([]byte)(event.Value), &pod)
		if err != nil {
			log.Println(err)
			return
		}
		if pod.Status.Phase == "stopping" {
			relativeURL := fmt.Sprintf("/api/v1/podInNode/%s/pods/%s/",
				pod.Spec.NodeName, pod.Metadata.Name)
			apiclient.Request(credential, relativeURL, nil, "DELETE")
			return
		}
		if pod.Status.Phase == "stopped" {
			relativeURL := fmt.Sprintf("/api/v1/pods/%s?remove=true",
				pod.Metadata.Name)
			apiclient.Request(credential, relativeURL, nil, "DELETE")
			return
		}
		if pod.Status.Phase != "" && pod.Status.Phase != "init" {
			return
		}
		node := selectNode()
		if node == nil {
			log.Println("no available node!")
			return
		}

		podIndex := podIndexOfNode[node.Metadata.Name] + 1
		if podIndex > maxPodIndexOfNode {
			log.Printf("node %s: too many pod\n", node.Metadata.Name)
		}
		podCIDR := node.Spec.PodCIDR
		podCIDRIP, podCIDRNet, err := net.ParseCIDR(podCIDR)
		if err != nil {
			log.Println(err)
			return
		}
		podIP := podCIDRIP.To4()
		if podIP == nil {
			log.Println(errors.New(fmt.Sprintf("%s: invalid IPv4 address", podCIDR)))
			return
		}
		if ones, _ := podCIDRNet.Mask.Size(); ones != 24 {
			log.Println(errors.New(fmt.Sprintf("%s: unrecognized pod CIDR network", podCIDR)))
			return
		}
		podIP[3] = byte(podIndex)
		podIndexOfNode[node.Metadata.Name] = podIndex

		pod.Status.Phase = "scheduled"
		pod.Spec.NodeName = node.Metadata.Name
		pod.Spec.IP = podIP.String()
		// TODO: fill other members
		apiclient.Request(credential,
			"/api/v1/pods/", pod, "PUT")
		apiclient.Request(credential,
			fmt.Sprintf("/api/v1/podInNode/%s/pods/", node.Metadata.Name),
			pod, "POST")
	}
}

// TODO: use strategy to choose a node
func selectNode() *ClusterResources.Node {
	nodeList := data.GetNodeList()
	if len(nodeList) == 0 {
		return nil
	}
	idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(nodeList))))
	return &nodeList[idx.Int64()]
}
