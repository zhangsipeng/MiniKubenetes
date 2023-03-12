package kubeproxy

import (
	"encoding/json"
	"errors"
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/data/ClusterResources"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"example/Minik8s/pkg/kubeproxy/iptables"
	"example/Minik8s/pkg/kubeproxy/vxlan"
	"example/Minik8s/pkg/kubeproxy/watchAPIServer"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

var ipt *iptables.IPTable

func initIPTables(reset bool) error {
	var err error
	ipt, err = iptables.New()
	if err != nil {
		return err
	}
	if reset {
		err = ipt.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

func initVxlan(runtimeConfig runtimedata.RuntimeConfig, reset bool) error {
	nodeId, ok := runtimeConfig.YamlConfig.Others["nodeId"]
	if !ok {
		return errors.New("no Others.nodeId found in config")
	}
	if _, err := exec.LookPath("ip"); err != nil {
		return err
	}
	if _, err := exec.LookPath("bridge"); err != nil {
		return err
	}
	var selfNode ClusterResources.Node
	var localVxlanIPCidr string
	err := json.Unmarshal(apiclient.Request(runtimeConfig,
		fmt.Sprintf("/api/v1/nodes/%s", nodeId), nil, "GET"),
		&selfNode)
	if err != nil {
		return err
	}
	if selfNode.Status.Phase == "init" {
		return errors.New("node hasn't been scheduled")
	}
	if reset {
		localVxlanIPCidr = selfNode.Spec.NodeVxlanCIDR
		if err := vxlan.DelVxlanIfExists(); err != nil {
			return err
		}
		if err := vxlan.InitVxlan(localVxlanIPCidr); err != nil {
			return err
		}
	}
	var nodes []ClusterResources.Node
	if err := json.Unmarshal(
		apiclient.Request(runtimeConfig, "/api/v1/nodes", nil, "GET"),
		&nodes); err != nil {
		return err
	}
	for _, peerNode := range nodes {
		if peerNode.Metadata.Name == selfNode.Metadata.Name {
			continue
		}
		if peerNode.Status.Phase == "init" { // no CIDR allocated
			continue
		}
		if err := vxlan.AddPeer(peerNode.Status.Addresses[0].Address,
			peerNode.Spec.NodeVxlanCIDR, peerNode.Spec.PodCIDR); err != nil {
			log.Printf("warning: %s", err.Error())
			continue
		}
	}
	return err
}

func initDNS(runtimeConfig runtimedata.RuntimeConfig, reset bool) (err error) {
	if reset {
		if err = exec.Command("resolvconf", "-d",
			fmt.Sprintf("%s.*", vxlan.VxlanIface)).Run(); err != nil {
			return
		}
	}
	cmd := exec.Command("resolvconf", "-a", fmt.Sprintf("%s.bind", vxlan.VxlanIface))
	cmd.Stdin = strings.NewReader(
		fmt.Sprintf("nameserver %s\nsearch pod.minik8s service.minik8s",
			runtimeConfig.YamlConfig.APIServerIP))
	err = cmd.Run()
	if err != nil {
		return
	}

	return
}

func StartService() {
	reset, info := getKubeProxyInitInfo()
	runtimeConfig := apiclient.InitRuntimeConfig(info, nil)
	err := initIPTables(reset)
	if err != nil {
		return
	}
	err = initVxlan(runtimeConfig, reset)
	if err != nil {
		panic(err)
	}
	err = initDNS(runtimeConfig, reset)
	go watchAPIServer.WatchNode(runtimeConfig)
	watchAPIServer.WatchService(ipt, runtimeConfig)
}

func getKubeProxyInitInfo() (bool, runtimedata.InitInfo) {
	reset := flag.Bool("reset", false, "是否重置iptables和vxlan")
	configFilename := flag.String("config", "default-config.yaml", "配置文件")
	flag.Parse()
	return *reset, runtimedata.InitInfo{
		Init:           false,
		ConfigFileName: *configFilename,
	}
}
