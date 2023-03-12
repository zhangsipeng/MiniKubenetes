package service

import (
	"encoding/json"
	"example/Minik8s/pkg/data/WorkloadResources"
	"example/Minik8s/pkg/kubeapiserver/service/const"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func createHorizontalPodAutoscaler(c *gin.Context) {
	var horizontalPodAutoscaler WorkloadResources.HorizontalPodAutoscaler
	err := c.BindJSON(&horizontalPodAutoscaler)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	horizontalPodAutoscalerJSON, _ := json.Marshal(horizontalPodAutoscaler)
	err = create(horizontalPodAutoscalerJSON, service_const.HorizontalPodAutoscalerPrefix, horizontalPodAutoscaler.Metadata.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, horizontalPodAutoscaler)
}

func updateHorizontalPodAutoscaler(c *gin.Context) {
	var horizontalPodAutoscaler WorkloadResources.HorizontalPodAutoscaler
	err := c.BindJSON(&horizontalPodAutoscaler)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	horizontalPodAutoscalerJSON, _ := json.Marshal(horizontalPodAutoscaler)
	err = put(horizontalPodAutoscalerJSON, service_const.HorizontalPodAutoscalerPrefix, horizontalPodAutoscaler.Metadata.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, horizontalPodAutoscaler)
}

func deleteHorizontalPodAutoscaler(c *gin.Context) {
	horizontalPodAutoscalerName := c.Param("name")
	if horizontalPodAutoscalerName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	err := delete(service_const.HorizontalPodAutoscalerPrefix, horizontalPodAutoscalerName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func listHorizontalPodAutoscaler(c *gin.Context) {
	ifWatch := c.DefaultQuery("watch", "false")

	// short HTTP connection, only list the horizontalPodAutoscalers
	if ifWatch == "false" {
		horizontalPodAutoscalerList := make([]WorkloadResources.HorizontalPodAutoscaler, 0)
		horizontalPodAutoscalerJSONList, err := list(service_const.HorizontalPodAutoscalerPrefix)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		for _, horizontalPodAutoscalerJSON := range horizontalPodAutoscalerJSONList {
			var horizontalPodAutoscaler WorkloadResources.HorizontalPodAutoscaler
			err = json.Unmarshal([]byte(horizontalPodAutoscalerJSON), &horizontalPodAutoscaler)
			if err != nil {
				continue
			}
			horizontalPodAutoscalerList = append(horizontalPodAutoscalerList, horizontalPodAutoscaler)
		}

		c.JSON(http.StatusOK, horizontalPodAutoscalerList)
		return
	}

	// long HTTP connection, create a watcher and start watching
	watchList(c, service_const.HorizontalPodAutoscalerPrefix)
}

func getHorizontalPodAutoscaler(c *gin.Context) {
	horizontalPodAutoscalerName := c.Param("name")
	if horizontalPodAutoscalerName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	horizontalPodAutoscalerJSON, err := get(service_const.HorizontalPodAutoscalerPrefix, horizontalPodAutoscalerName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if horizontalPodAutoscalerJSON != "" {
		var horizontalPodAutoscaler WorkloadResources.HorizontalPodAutoscaler
		err = json.Unmarshal([]byte(horizontalPodAutoscalerJSON), &horizontalPodAutoscaler)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, horizontalPodAutoscaler)
		return
	}

	c.JSON(http.StatusOK, nil)
}
