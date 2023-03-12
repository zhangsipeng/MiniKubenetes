package service

import (
	"encoding/json"
	"example/Minik8s/pkg/data/Serverless"
	"example/Minik8s/pkg/kubeapiserver/service/const"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func createServerlessService(c *gin.Context) {
	var serverlessService Serverless.Service
	err := c.BindJSON(&serverlessService)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	serverlessService.Status.Phase = "init"
	serverlessServiceJSON, _ := json.Marshal(serverlessService)
	err = create(serverlessServiceJSON, service_const.ServerlessServicePrefix, serverlessService.Metadata.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, serverlessService)
}

func updateServerlessService(c *gin.Context) {
	var serverlessService Serverless.Service
	err := c.BindJSON(&serverlessService)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	serverlessServiceJSON, _ := json.Marshal(serverlessService)
	err = put(serverlessServiceJSON, service_const.ServerlessServicePrefix, serverlessService.Metadata.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, serverlessService)
}

func deleteServerlessService(c *gin.Context) {
	serverlessServiceName := c.Param("name")
	if serverlessServiceName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	ifRemove := c.DefaultQuery("remove", "false")
	if ifRemove == "false" {
		serviceJSON, err := get(service_const.ServerlessServicePrefix, serverlessServiceName)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		if serviceJSON == "" {
			log.Println(service_const.NotExistError)
			c.JSON(http.StatusBadRequest, service_const.NotExistError)
			return
		}

		var service Serverless.Service
		err = json.Unmarshal([]byte(serviceJSON), &service)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		service.Status.Phase = "removing"
		serviceJSON_, _ := json.Marshal(service)
		err = put(serviceJSON_, service_const.ServerlessServicePrefix, service.Metadata.Name)

		c.JSON(http.StatusOK, nil)
		return
	}

	err := delete(service_const.ServerlessServicePrefix, serverlessServiceName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func listServerlessService(c *gin.Context) {
	ifWatch := c.DefaultQuery("watch", "false")

	// short HTTP connection, only list the serverlessServices
	if ifWatch == "false" {
		serverlessServiceList := make([]Serverless.Service, 0)
		serverlessServiceJSONList, err := list(service_const.ServerlessServicePrefix)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		for _, serverlessServiceJSON := range serverlessServiceJSONList {
			var serverlessService Serverless.Service
			err = json.Unmarshal([]byte(serverlessServiceJSON), &serverlessService)
			if err != nil {
				continue
			}
			serverlessServiceList = append(serverlessServiceList, serverlessService)
		}

		c.JSON(http.StatusOK, serverlessServiceList)
		return
	}

	// long HTTP connection, create a watcher and start watching
	watchList(c, service_const.ServerlessServicePrefix)
}

func getServerlessService(c *gin.Context) {
	serverlessServiceName := c.Param("name")
	if serverlessServiceName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	serverlessServiceJSON, err := get(service_const.ServerlessServicePrefix, serverlessServiceName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if serverlessServiceJSON != "" {
		var serverlessService Serverless.Service
		err = json.Unmarshal([]byte(serverlessServiceJSON), &serverlessService)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, serverlessService)
		return
	}

	c.JSON(http.StatusOK, nil)
}
