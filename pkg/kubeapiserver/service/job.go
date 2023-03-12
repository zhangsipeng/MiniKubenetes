package service

import (
	"encoding/json"
	"example/Minik8s/pkg/data/WorkloadResources"
	service_const "example/Minik8s/pkg/kubeapiserver/service/const"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func createJob(c *gin.Context) {
	var job WorkloadResources.GPUJob
	err := c.BindJSON(&job)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	job.Phase = "init"
	jobJSON, _ := json.Marshal(job)
	err = create(jobJSON, service_const.JobPrefix, job.Metadata.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, job)
}

func createJobInNode(c *gin.Context) {
	nodeName := c.Param("name")
	if nodeName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	var job WorkloadResources.GPUJob
	err := c.BindJSON(&job)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	prefix := service_const.JobInNodePrefix + nodeName + "/"
	err = create([]byte{}, prefix, job.Metadata.Name) // only a key, no value
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, job)
}

func updateJob(c *gin.Context) {
	var job WorkloadResources.GPUJob
	err := c.BindJSON(&job)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	jobJSON, _ := json.Marshal(job)
	err = put(jobJSON, service_const.JobPrefix, job.Metadata.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, job)
}

func deleteJob(c *gin.Context) {
	jobName := c.Param("name")
	if jobName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	ifRemove := c.DefaultQuery("remove", "false")
	if ifRemove == "false" {
		jobJSON, err := get(service_const.JobPrefix, jobName)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		if jobJSON == "" {
			log.Println(service_const.NotExistError)
			c.JSON(http.StatusBadRequest, service_const.NotExistError)
			return
		}

		var job WorkloadResources.GPUJob
		err = json.Unmarshal([]byte(jobJSON), &job)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		job.Phase = "stopping"
		jobJSON_, _ := json.Marshal(job)
		err = put(jobJSON_, service_const.JobPrefix, job.Metadata.Name)

		c.JSON(http.StatusOK, nil)
		return
	}

	err := delete(service_const.JobPrefix, jobName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func deleteJobInNode(c *gin.Context) {
	nodeName := c.Param("nodeName")
	jobName := c.Param("jobName")
	if nodeName == "" || jobName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	err := delete(service_const.JobInNodePrefix+nodeName+"/", jobName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func listJob(c *gin.Context) {
	ifWatch := c.DefaultQuery("watch", "false")

	// short HTTP connection, only list the jobs
	if ifWatch == "false" {
		jobList := make([]WorkloadResources.GPUJob, 0)
		jobJSONList, err := list(service_const.JobPrefix)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		for _, jobJSON := range jobJSONList {
			var job WorkloadResources.GPUJob
			err = json.Unmarshal([]byte(jobJSON), &job)
			if err != nil {
				continue
			}
			jobList = append(jobList, job)
		}

		c.JSON(http.StatusOK, jobList)
		return
	}

	// long HTTP connection, create a watcher and start watching
	watchList(c, service_const.JobPrefix)
}

func listJobInNode(c *gin.Context) {
	nodeName := c.Param("name")
	if nodeName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	ifWatch := c.DefaultQuery("watch", "false")

	// short HTTP connection, only list the jobs
	if ifWatch == "false" {
		jobList := make([]WorkloadResources.GPUJob, 0)
		jobJSONList, err := list(service_const.JobInNodePrefix + nodeName + "/")
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		for _, jobJSON := range jobJSONList {
			var job WorkloadResources.GPUJob
			err = json.Unmarshal([]byte(jobJSON), &job)
			if err != nil {
				continue
			}
			jobList = append(jobList, job)
		}

		c.JSON(http.StatusOK, jobList)
		return
	}

	// long HTTP connection, create a watcher and start watching
	watchList(c, service_const.JobInNodePrefix+nodeName+"/")
}

func getJob(c *gin.Context) {
	jobName := c.Param("name")
	if jobName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	jobJSON, err := get(service_const.JobPrefix, jobName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if jobJSON != "" {
		var job WorkloadResources.GPUJob
		err = json.Unmarshal([]byte(jobJSON), &job)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, job)
		return
	}

	c.JSON(http.StatusOK, nil)
}
