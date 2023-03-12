package server

import (
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/serverless/data"
	"log"
	"time"
)

func startMonitorService(service *serverlessServiceInfo) {
	for {
		if service.ifStop == true {
			log.Println("stop monitor " + service.service.Metadata.Name)
			break
		}
		time.Sleep(30 * time.Second)
		podNum := len(service.pods)
		reqNum := service.requestNum
		service.requestNum = 0
		expectedNum := (reqNum + 9) / 10
		if podNum < expectedNum {
			more := expectedNum - podNum
			for i := 1; i <= more; i++ {
				newPod := new(podInfo)
				newPod.pod = createPod(service)
				newPod.activeReq = 0
				apiclient.Request(data.GetCredential(), "/api/v1/pods/", newPod.pod, "POST")
				service.pods = append(service.pods, newPod)
			}
		} else if podNum > expectedNum {
			deletePodList := service.pods[:podNum-expectedNum]
			service.pods = service.pods[podNum-expectedNum:]
			go deletePods(deletePodList)
		}
	}
}

func deletePods(podList []*podInfo) {
	for _, pod := range podList {
		go deletePod(pod)
	}
}

func deletePod(pod *podInfo) {
	waitPodFree(pod)
	log.Println("delete pod " + pod.pod.Metadata.Name)
	apiclient.Request(data.GetCredential(), "/api/v1/pods/"+pod.pod.Metadata.Name, nil, "DELETE")
}

func waitPodFree(pod *podInfo) {
	for {
		if pod.activeReq == 0 {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
}
