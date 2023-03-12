package service

import (
	"encoding/json"
	"example/Minik8s/pkg/data/ServiceResources"
	service_const "example/Minik8s/pkg/kubeapiserver/service/const"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func createService(c *gin.Context) {
	var service ServiceResources.Service
	err := c.BindJSON(&service)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	serviceJSON, _ := json.Marshal(service)
	err = create(serviceJSON, service_const.ServicePrefix, service.Metadata.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, service)
}

func updateService(c *gin.Context) {
	var service ServiceResources.Service
	err := c.BindJSON(&service)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	serviceJSON, _ := json.Marshal(service)
	err = put(serviceJSON, service_const.ServicePrefix, service.Metadata.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, service)
}

func deleteService(c *gin.Context) {
	serviceName := c.Param("name")
	if serviceName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	err := delete(service_const.ServicePrefix, serviceName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func listService(c *gin.Context) {
	ifWatch := c.DefaultQuery("watch", "false")

	// short HTTP connection, only list the services
	if ifWatch == "false" {
		serviceList := make([]ServiceResources.Service, 0)
		serviceJSONList, err := list(service_const.ServicePrefix)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		for _, serviceJSON := range serviceJSONList {
			var service ServiceResources.Service
			err = json.Unmarshal([]byte(serviceJSON), &service)
			if err != nil {
				continue
			}
			serviceList = append(serviceList, service)
		}

		c.JSON(http.StatusOK, serviceList)
		return
	}

	// long HTTP connection, create a watcher and start watching
	watchList(c, service_const.ServicePrefix)
}

func getService(c *gin.Context) {
	serviceName := c.Param("name")
	if serviceName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	serviceJSON, err := get(service_const.ServicePrefix, serviceName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if serviceJSON != "" {
		var service ServiceResources.Service
		err = json.Unmarshal([]byte(serviceJSON), &service)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, service)
		return
	}

	c.JSON(http.StatusOK, nil)
}
