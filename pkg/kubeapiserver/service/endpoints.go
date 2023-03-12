package service

import (
	"encoding/json"
	"example/Minik8s/pkg/data/ServiceResources"
	"example/Minik8s/pkg/kubeapiserver/service/const"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func createEndpoints(c *gin.Context) {
	var endpoints ServiceResources.Endpoints
	err := c.BindJSON(&endpoints)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	endpointsJSON, _ := json.Marshal(endpoints)
	err = create(endpointsJSON, service_const.EndpointsPrefix, endpoints.Metadata.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, endpoints)
}

func updateEndpoints(c *gin.Context) {
	var endpoints ServiceResources.Endpoints
	err := c.BindJSON(&endpoints)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	endpointsJSON, _ := json.Marshal(endpoints)
	err = put(endpointsJSON, service_const.EndpointsPrefix, endpoints.Metadata.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, endpoints)
}

func deleteEndpoints(c *gin.Context) {
	endpointsName := c.Param("name")
	if endpointsName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	err := delete(service_const.EndpointsPrefix, endpointsName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func listEndpoints(c *gin.Context) {
	ifWatch := c.DefaultQuery("watch", "false")

	// short HTTP connection, only list the endpointss
	if ifWatch == "false" {
		endpointsList := make([]ServiceResources.Endpoints, 0)
		endpointsJSONList, err := list(service_const.EndpointsPrefix)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		for _, endpointsJSON := range endpointsJSONList {
			var endpoints ServiceResources.Endpoints
			err = json.Unmarshal([]byte(endpointsJSON), &endpoints)
			if err != nil {
				continue
			}
			endpointsList = append(endpointsList, endpoints)
		}

		c.JSON(http.StatusOK, endpointsList)
		return
	}

	// long HTTP connection, create a watcher and start watching
	watchList(c, service_const.EndpointsPrefix)
}

func getEndpoints(c *gin.Context) {
	endpointsName := c.Param("name")
	if endpointsName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	ifWatch := c.DefaultQuery("watch", "false")
	if ifWatch == "false" {
		endpointsJSON, err := get(service_const.EndpointsPrefix, endpointsName)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		if endpointsJSON != "" {
			var endpoints ServiceResources.Endpoints
			err = json.Unmarshal([]byte(endpointsJSON), &endpoints)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			c.JSON(http.StatusOK, endpoints)
			return
		}

		c.JSON(http.StatusOK, nil)
		return
	}

	watchObject(c, service_const.EndpointsPrefix, endpointsName)
}
