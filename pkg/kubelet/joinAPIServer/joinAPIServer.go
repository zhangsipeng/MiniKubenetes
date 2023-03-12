package joinAPIServer

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/data/ClusterResources"
	"example/Minik8s/pkg/data/ObjectMeta"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"fmt"
	"time"
)

const NODE_NAME_BYTE = 30

func JoinAPIServer(config runtimedata.RuntimeConfig) {
	nodeId, exist := config.YamlConfig.Others["nodeId"]
	if !exist {
		rand_byte := make([]byte, NODE_NAME_BYTE)
		_, err := rand.Read(rand_byte)
		if err != nil {
			panic(err)
		}
		nodeId = base64.RawURLEncoding.EncodeToString(rand_byte)
	}
	nodeInfo := ClusterResources.Node{
		// TODO: missing field
		Kind: "node",
		Metadata: ObjectMeta.ObjectMeta{
			Name: nodeId,
		},
		Status: ClusterResources.NodeStatus{
			Addresses: []ClusterResources.NodeAddress{{
				Address: config.YamlConfig.ClientIP,
			}},
			Phase: "init",
		},
	}
	apiclient.Request(config, "/api/v1/nodes", nodeInfo, "POST")
}

func WaitForScheduler(runtimeConfig runtimedata.RuntimeConfig) (*ClusterResources.Node, error) {
	nodeId, exist := runtimeConfig.YamlConfig.Others["nodeId"]
	if !exist {
		return nil, errors.New("no Others.nodeId found in runtimeConfig")
	}
	var node ClusterResources.Node
	for {
		err := json.Unmarshal(apiclient.Request(runtimeConfig,
			fmt.Sprintf("/api/v1/nodes/%s", nodeId), nil, "GET"),
			&node)
		if err != nil {
			return nil, err
		}
		if node.Status.Phase != "init" {
			break
		}
		time.Sleep(5 * time.Second)
	}
	return &node, nil
}
