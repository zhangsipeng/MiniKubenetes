package service

import (
	"encoding/json"
	"example/Minik8s/pkg/data/WorkloadResources"
	"example/Minik8s/pkg/kubeapiserver/service/const"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func createReplicaSet(c *gin.Context) {
	var replicaSet WorkloadResources.ReplicaSet
	err := c.BindJSON(&replicaSet)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	replicaSetJSON, _ := json.Marshal(replicaSet)
	err = create(replicaSetJSON, service_const.ReplicaSetPrefix, replicaSet.Metadata.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, replicaSet)
}

func updateReplicaSet(c *gin.Context) {
	var replicaSet WorkloadResources.ReplicaSet
	err := c.BindJSON(&replicaSet)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	replicaSetJSON, _ := json.Marshal(replicaSet)
	err = put(replicaSetJSON, service_const.ReplicaSetPrefix, replicaSet.Metadata.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, replicaSet)
}

func deleteReplicaSet(c *gin.Context) {
	replicaSetName := c.Param("name")
	if replicaSetName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	err := delete(service_const.ReplicaSetPrefix, replicaSetName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func listReplicaSet(c *gin.Context) {
	ifWatch := c.DefaultQuery("watch", "false")

	// short HTTP connection, only list the replicaSets
	if ifWatch == "false" {
		replicaSetList := make([]WorkloadResources.ReplicaSet, 0)
		replicaSetJSONList, err := list(service_const.ReplicaSetPrefix)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		for _, replicaSetJSON := range replicaSetJSONList {
			var replicaSet WorkloadResources.ReplicaSet
			err = json.Unmarshal([]byte(replicaSetJSON), &replicaSet)
			if err != nil {
				continue
			}
			replicaSetList = append(replicaSetList, replicaSet)
		}

		c.JSON(http.StatusOK, replicaSetList)
		return
	}

	// long HTTP connection, create a watcher and start watching
	watchList(c, service_const.ReplicaSetPrefix)
}

func getReplicaSet(c *gin.Context) {
	replicaSetName := c.Param("name")
	if replicaSetName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	replicaSetJSON, err := get(service_const.ReplicaSetPrefix, replicaSetName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if replicaSetJSON != "" {
		var replicaSet WorkloadResources.ReplicaSet
		err = json.Unmarshal([]byte(replicaSetJSON), &replicaSet)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, replicaSet)
		return
	}

	c.JSON(http.StatusOK, nil)
}
