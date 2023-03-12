package service

import (
	"encoding/json"
	"example/Minik8s/pkg/data/WorkloadResources"
	service_const "example/Minik8s/pkg/kubeapiserver/service/const"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createPod(c *gin.Context) {
	var pod WorkloadResources.Pod
	err := c.BindJSON(&pod)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	podJSON, _ := json.Marshal(pod)
	err = create(podJSON, service_const.PodPrefix, pod.Metadata.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, pod)
}

func createPodInNode(c *gin.Context) {
	nodeName := c.Param("name")
	if nodeName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	var pod WorkloadResources.Pod
	err := c.BindJSON(&pod)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	prefix := service_const.PodInNodePrefix + nodeName + "/"
	err = create([]byte{}, prefix, pod.Metadata.Name) // only a key, no value
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, pod)
}

func updatePod(c *gin.Context) {
	var pod WorkloadResources.Pod
	err := c.BindJSON(&pod)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	podJSON, _ := json.Marshal(pod)
	err = put(podJSON, service_const.PodPrefix, pod.Metadata.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, pod)
}

func deletePod(c *gin.Context) {
	podName := c.Param("name")
	if podName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	ifRemove := c.DefaultQuery("remove", "false")
	if ifRemove == "false" {
		podJSON, err := get(service_const.PodPrefix, podName)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		if podJSON == "" {
			log.Println(service_const.NotExistError)
			c.JSON(http.StatusBadRequest, service_const.NotExistError)
			return
		}

		var pod WorkloadResources.Pod
		err = json.Unmarshal([]byte(podJSON), &pod)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		pod.Status.Phase = "stopping"
		podJSON_, _ := json.Marshal(pod)
		err = put(podJSON_, service_const.PodPrefix, pod.Metadata.Name)

		c.JSON(http.StatusOK, nil)
		return
	}

	err := delete(service_const.PodPrefix, podName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func deletePodInNode(c *gin.Context) {
	nodeName := c.Param("nodeName")
	podName := c.Param("podName")
	if nodeName == "" || podName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	err := delete(service_const.PodInNodePrefix+nodeName+"/", podName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func listPod(c *gin.Context) {
	ifWatch := c.DefaultQuery("watch", "false")

	// short HTTP connection, only list the pods
	if ifWatch == "false" {
		podList := make([]WorkloadResources.Pod, 0)
		podJSONList, err := list(service_const.PodPrefix)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		for _, podJSON := range podJSONList {
			var pod WorkloadResources.Pod
			err = json.Unmarshal([]byte(podJSON), &pod)
			if err != nil {
				continue
			}
			podList = append(podList, pod)
		}

		c.JSON(http.StatusOK, podList)
		return
	}

	// long HTTP connection, create a watcher and start watching
	watchList(c, service_const.PodPrefix)
}

func listPodInNode(c *gin.Context) {
	nodeName := c.Param("name")
	if nodeName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	ifWatch := c.DefaultQuery("watch", "false")

	// short HTTP connection, only list the pods
	if ifWatch == "false" {
		podList := make([]WorkloadResources.Pod, 0)
		podJSONList, err := list(service_const.PodInNodePrefix + nodeName + "/")
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		for _, podJSON := range podJSONList {
			var pod WorkloadResources.Pod
			err = json.Unmarshal([]byte(podJSON), &pod)
			if err != nil {
				continue
			}
			podList = append(podList, pod)
		}

		c.JSON(http.StatusOK, podList)
		return
	}

	// long HTTP connection, create a watcher and start watching
	watchList(c, service_const.PodInNodePrefix+nodeName+"/")
}

func getPod(c *gin.Context) {
	podName := c.Param("name")
	if podName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	podJSON, err := get(service_const.PodPrefix, podName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if podJSON != "" {
		var pod WorkloadResources.Pod
		err = json.Unmarshal([]byte(podJSON), &pod)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, pod)
		return
	}

	c.JSON(http.StatusOK, nil)
}
