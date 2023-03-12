package reportNodeStatus

import (
	"encoding/json"
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/data/ClusterResources"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"fmt"
	"time"
)

func HeartBeat(runtimeConfig runtimedata.RuntimeConfig) {
	nodeId, ok := runtimeConfig.YamlConfig.Others["nodeId"]
	if !ok {
		panic("no YamlConfig.Others.nodeId in runtimeConfig")
	}
	for {
		var node ClusterResources.Node
		if err := json.Unmarshal(apiclient.Request(runtimeConfig,
			fmt.Sprintf("/api/v1/nodes/%s", nodeId), nil, "GET"),
			&node); err != nil {
			panic(err)
		}
		node.Status.Conditions = []ClusterResources.NodeCondition{
			{
				LastHeartBeatTime: time.Now(),
			},
		}
		apiclient.Request(runtimeConfig, "/api/v1/nodes/", node, "PUT")
		time.Sleep(5 * time.Second)
	}
}
