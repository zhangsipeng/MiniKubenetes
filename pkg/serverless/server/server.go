package server

import (
	"crypto/rand"
	"encoding/json"
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/data/ObjectMeta"
	"example/Minik8s/pkg/data/Serverless"
	"example/Minik8s/pkg/data/WorkloadResources"
	"example/Minik8s/pkg/serverless/data"
	"example/Minik8s/pkg/serverless/docker"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const podPort = 8080

var client *http.Client

type podInfo struct {
	pod       WorkloadResources.Pod
	activeReq int
}

type serverlessServiceInfo struct {
	service    Serverless.Service
	requestNum int
	pods       []*podInfo
	podIndex   int
	imageName  string
	ifStop     bool
}

var serviceMap map[string]*serverlessServiceInfo

func InitServer() {
	client = new(http.Client)
	serviceMap = make(map[string]*serverlessServiceInfo)
	DAGMap = make(map[string]*Serverless.DAG)
	AddService("GET", "/service/:name/", handleServiceRequest)
	AddService("GET", "/DAG/:name/", handleDAGRequest)
}

func StartServerlessService(serverlessService Serverless.Service) {
	imageName, err := docker.GenerateDockerImage(serverlessService)
	if err != nil {
		panic(err)
	}
	serviceName := serverlessService.Metadata.Name
	_, ok := serviceMap[serviceName]
	if ok {
		log.Println("service " + serviceName + " already exists!")
		return
	}

	newService := new(serverlessServiceInfo)

	newService.service.ApiVersion = "v1"
	newService.service.Kind = "service"
	newService.service.Metadata = serverlessService.Metadata
	newService.service.Spec.GitUrl = serverlessService.Spec.GitUrl
	newService.service.Spec.MaxReplicas = serverlessService.Spec.MaxReplicas
	newService.service.Spec.Input = serverlessService.Spec.Input
	newService.service.Spec.Output = serverlessService.Spec.Output
	newService.service.Status.Phase = "running"
	newService.service.Status.Replicas = 0

	newService.requestNum = 0
	newService.pods = make([]*podInfo, 0)
	newService.imageName = imageName

	newService.ifStop = false

	serviceMap[serviceName] = newService

	go startMonitorService(newService)
	apiclient.Request(data.GetCredential(), "/api/v1/serverless/", newService.service, "PUT")
	log.Println("add service " + newService.service.Metadata.Name)
}

func UpdateServerlessService(serverlessService Serverless.Service) {
	DeleteServerlessService(serverlessService, true)
	StartServerlessService(serverlessService)
}

func DeleteServerlessService(serverlessService Serverless.Service, ifUpdate bool) {
	serviceName := serverlessService.Metadata.Name
	targetService, ok := serviceMap[serviceName]
	if !ok {
		log.Println("no such service")
		return
	}
	targetService.ifStop = true
	deletePods(targetService.pods)
	delete(serviceMap, serviceName)
	if !ifUpdate {
		relativeURL := fmt.Sprintf("/api/v1/serverless/%s/?remove=true",
			serviceName)
		apiclient.Request(data.GetCredential(), relativeURL, nil, "DELETE")
	}
	log.Println("delete service " + serviceName)
}

func handleServiceRequest(c *gin.Context) {
	serviceName := c.Param("name")
	if serviceName == "" {
		log.Println("empty service name")
		c.JSON(http.StatusBadRequest, "empty service name")
		return
	}

	serverlessService, ok := serviceMap[serviceName]
	if !ok {
		log.Println("no such service")
		c.JSON(http.StatusBadRequest, "no such service")
		return
	}

	query := make(map[string]string)
	for _, key := range serverlessService.service.Spec.Input {
		val := c.DefaultQuery(key, "")
		query[key] = val
	}

	response, err := forwardRequest(serverlessService, query)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	defer response.Body.Close()

	resBody, _ := ioutil.ReadAll(response.Body)
	c.Data(response.StatusCode, response.Header.Get("content-type"), resBody)
}

func forwardRequest(serverlessService *serverlessServiceInfo, query map[string]string) (*http.Response, error) {
	serverlessService.requestNum = serverlessService.requestNum + 1
	var selectPod *podInfo
	if len(serverlessService.pods) == 0 {
		// create a new pod to handle the request
		newPod := new(podInfo)
		newPod.pod = createPod(serverlessService)
		newPod.activeReq = 0
		apiclient.Request(data.GetCredential(), "/api/v1/pods/", newPod.pod, "POST")
		serverlessService.pods = append(serverlessService.pods, newPod)
		selectPod = newPod
	} else {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(serverlessService.pods))))
		selectPod = serverlessService.pods[idx.Int64()]
	}

	selectPod.activeReq = selectPod.activeReq + 1
	// wait until the pod is ready
	log.Println("new request to " + selectPod.pod.Metadata.Name)
	ip, err := waitPodStart(selectPod.pod.Metadata.Name)
	if err != nil {
		return nil, err
	}

	targetUrl := fmt.Sprintf("http://%s:%d", ip, podPort)
	request, _ := http.NewRequest("GET", targetUrl, nil)

	params := make(url.Values)
	for key, val := range query {
		params.Add(key, val)
	}
	request.URL.RawQuery = params.Encode()

	var response *http.Response
	for i := 1; i <= 3; i++ {
		response, err = client.Do(request)
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	selectPod.activeReq = selectPod.activeReq - 1
	if err != nil {
		return nil, err
	}
	return response, nil
}

func createPod(service *serverlessServiceInfo) WorkloadResources.Pod {
	service.podIndex = service.podIndex + 1
	newPod := WorkloadResources.Pod{
		ApiVersion: "v1",
		Kind:       "pod",
		Metadata: ObjectMeta.ObjectMeta{
			Name: "serverless-" + service.service.Metadata.Name + "-pod" + strconv.Itoa(service.podIndex),
		},
		Spec: WorkloadResources.PodSpec{
			Container: []WorkloadResources.Container{
				{
					Name:  "serverless-service",
					Image: service.imageName,
					Ports: []WorkloadResources.ContainerPort{
						{
							Protocal:      "tcp",
							ContainerPort: podPort,
						},
					},
				},
			},
		},
	}
	log.Println("create pod " + newPod.Metadata.Name)
	return newPod
}

func waitPodStart(podName string) (string, error) {
	var pod WorkloadResources.Pod
	for {
		res := apiclient.Request(data.GetCredential(),
			fmt.Sprintf("/api/v1/pods/%s/", podName), nil, "GET")
		err := json.Unmarshal(res, &pod)
		if err != nil {
			return "", err
		}
		if pod.Status.Phase == "running" {
			return pod.Spec.IP, nil
		}
		time.Sleep(500 * time.Millisecond)
	}
}
