package data

import (
	"example/Minik8s/pkg/data/ClusterResources"
	"fmt"
	"log"
)

var nodeMap = make(map[string]ClusterResources.Node)
var nodeList = make([]ClusterResources.Node, 0)

var nextAvailNodeIdx = 1

const maxNodeIdx = 254

func AddNode(node *ClusterResources.Node) {
	if nextAvailNodeIdx > maxNodeIdx {
		log.Printf("error when adding node %s: too many node", node.Metadata.Name)
		return
	}
	node.Status.Phase = "scheduled"
	node.Spec.NodeVxlanCIDR = fmt.Sprintf("10.37.4.%d/24", nextAvailNodeIdx)
	node.Spec.PodCIDR = fmt.Sprintf("172.37.%d.0/24", nextAvailNodeIdx)
	nextAvailNodeIdx += 1
	if _, ok := nodeMap[node.Metadata.Name]; !ok {
		nodeMap[node.Metadata.Name] = *node
		nodeList = append(nodeList, *node)
	}
	log.Println("add node " + node.Metadata.Name)
}

func GetNodeList() []ClusterResources.Node {
	return nodeList
}
