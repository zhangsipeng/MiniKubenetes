package main

import (
	"bytes"
	"encoding/json"
	"example/Minik8s/pkg/data/ObjectMeta"
	"example/Minik8s/pkg/data/WorkloadResources"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func main() {
	testWatch()
}

func testWatch() {
	targetUrl := "http://localhost:8080/api/v1/pods"
	request, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		log.Fatalln(err)
	}

	params := make(url.Values)
	params.Add("watch", "true")

	request.URL.RawQuery = params.Encode()

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	buf := make([]byte, 4096) // any non zero value will do, try '1'.
	for {
		n, err := response.Body.Read(buf)
		if n == 0 && err != nil { // simplified
			break
		}
	}
}

func testGet() {
	targetUrl := "http://localhost:8080/api/v1/pods"
	request, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		log.Fatalln(err)
	}

	params := make(url.Values)
	params.Add("watch", "false")

	request.URL.RawQuery = params.Encode()
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	var podList []WorkloadResources.Pod

	s, _ := ioutil.ReadAll(response.Body)

	err = json.Unmarshal(s, &podList)
	if err != nil {
		return
	}
}

func testCreate() {
	targetUrl := "http://localhost:8080/api/v1/pods"
	pod := WorkloadResources.Pod{
		ApiVersion: "v1",
		Kind:       "pod",
		Metadata: ObjectMeta.ObjectMeta{
			Name: "nginx",
		},
		Spec: WorkloadResources.PodSpec{
			Container: []WorkloadResources.Container{
				{
					Name:  "nginx",
					Image: "nginx:1.14.2",
					Ports: []WorkloadResources.ContainerPort{
						{
							ContainerPort: 80,
						},
					},
				},
			},
		},
	}

	podJSON, err := json.Marshal(pod)
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest("POST", targetUrl, bytes.NewReader(podJSON))

	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(response.Body)
}
