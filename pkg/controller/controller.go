package controller

import (
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/controller/watchAPIServer"
	"example/Minik8s/pkg/data/ClusterResources"
	"example/Minik8s/pkg/data/WorkloadResources"
	"fmt"
	"log"
	"time"
)

func StartService() {
	info := apiclient.GetInitInfo()
	runtimeConfig := apiclient.InitRuntimeConfig(info, nil)
	go watchAPIServer.WatchService(runtimeConfig)
	type any = interface{}
	for {

		//发送请求，获得deployment信息
		var deploymentlist []WorkloadResources.Deployment
		response := apiclient.Request(runtimeConfig, "/api/v1/deployment", nil, "GET")
		str_response := string(response[:])
		if str_response != "null" {

			err := json.Unmarshal(response, &deploymentlist)
			if err != nil {
				panic("Fatel:decode deploymentlist json failed \n")
			}
			//对每个deployment查看
			log.Printf("receive deployment list: %s\n", string(response))
			for _, item := range deploymentlist {
				livingPod := []string{}
				for _, podName := range item.Status.PodName {
					Url := "/api/v1/pods/" + podName
					response := apiclient.Request(runtimeConfig, Url, nil, "GET")
					str_res := string(response[:])
					if str_res == "null" {
						log.Printf("pod %s died, response null\n", podName)
						continue // dead
					}
					var targetPod WorkloadResources.Pod
					err := json.Unmarshal(response, &targetPod)
					if err != nil {
						panic("Fatel:decode pod json failed \n")
					}
					if targetPod.Status.Phase == "stopped" {
						log.Printf("pod %s died, phase stopped\n", podName)
						continue //dead
					}
					//pod还在运行，看一下node是否存活
					Url = "/api/v1/nodes/" + targetPod.Spec.NodeName
					response = apiclient.Request(runtimeConfig, Url, nil, "GET")
					var targetNode ClusterResources.Node
					err = json.Unmarshal(response, &targetNode)
					if err != nil {
						//TODO:may be fix something?
						panic("Fatel:decode Node json failed \n")
					}
					if len(targetNode.Status.Conditions) != 0 {
						nowTime := time.Now()
						if (nowTime.Sub(targetNode.Status.Conditions[0].LastHeartBeatTime)).Seconds() > 60 {
							log.Printf("pod %s died, node %s inactive\n", podName, targetPod.Spec.NodeName)
							continue // dead
						}
					}
					livingPod = append(livingPod, podName)
				}
				livingPodCnt := len(livingPod)

				//创建max(0, spec.replicas - livintPodCnt)数目个的pod
				podToCreate := item.Spec.Replicas - int32(livingPodCnt)
				if podToCreate < 0 {
					podToCreate = 0
				}
				log.Printf("deployment %s need %d pods, has %d pods, %d to start\n",
					item.Metadata.Name,
					item.Spec.Replicas, livingPodCnt, podToCreate)

				for K := 0; K < int(podToCreate); K++ {
					var templatePod WorkloadResources.Pod
					templatePod.ApiVersion = item.ApiVersion
					//TODO:lower upper case matter?
					templatePod.Kind = "pod"
					templatePod.Metadata.Labels = item.Spec.Template.Metadata.Labels

					rand_byte := make([]byte, 4)
					_, err := rand.Read(rand_byte)
					if err != nil {
						panic(err)
					}
					randomName :=
						base32.HexEncoding.WithPadding(base32.NoPadding).EncodeToString(rand_byte)

					templatePod.Metadata.Name = apiclient.MangleName("deployment", item.Metadata.Name, randomName)
					templatePod.Spec = item.Spec.Template.Spec
					response = apiclient.Request(runtimeConfig, "/api/v1/pods", templatePod, "POST")
					fmt.Printf("pod create response :\n%s\n", response)
					livingPod = append(livingPod, templatePod.Metadata.Name)
				}
				item.Status.Replicas = int32(len(livingPod))
				item.Status.PodName = livingPod
				apiclient.Request(runtimeConfig, "/api/v1/deployment/", item, "PUT")

			}
		}
		time.Sleep(20 * time.Second)
	}
}
