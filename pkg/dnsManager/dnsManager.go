package dnsManager

import (
	"encoding/json"
	"example/Minik8s/pkg/apiclient"
	runtimedata "example/Minik8s/pkg/data/RuntimeData"
	"example/Minik8s/pkg/data/WorkloadResources"
	"example/Minik8s/pkg/dnsManager/dnsOp"
	"example/Minik8s/pkg/dnsManager/watch"
	"flag"
)

func initDns(runtimeConfig runtimedata.RuntimeConfig) error {
	var pods []WorkloadResources.Pod
	if err := json.Unmarshal(
		apiclient.Request(runtimeConfig, "/api/v1/pods/", nil, "GET"),
		&pods); err != nil {
		return err
	}
	for _, pod := range pods {
		dnsOp.AddPod(pod)
	}
	return nil
}

func StartService() {
	ddnsKeyFilename, info := getDnsManagerInitInfo()
	runtimeConfig := apiclient.InitRuntimeConfig(info, nil)
	dnsOp.InitDnsOp(runtimeConfig.YamlConfig.APIServerIP, ddnsKeyFilename)
	go watch.WatchPod(runtimeConfig)
	watch.WatchService(runtimeConfig)
}

func getDnsManagerInitInfo() (*string, runtimedata.InitInfo) {
	ddnsKeyFilename := flag.String("ddns-key", "", "DDNS更新使用的密钥")
	info := apiclient.GetInitInfo()
	if *ddnsKeyFilename == "" {
		ddnsKeyFilename = nil
	}
	return ddnsKeyFilename, info
}
