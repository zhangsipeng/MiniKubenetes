package service

import (
	"encoding/json"
	"example/Minik8s/pkg/data/Serverless"
	"example/Minik8s/pkg/kubeapiserver/service/const"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func createServerlessDAG(c *gin.Context) {
	var serverlessDAG Serverless.DAG
	err := c.BindJSON(&serverlessDAG)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	serverlessDAG.Status.Phase = "init"
	serverlessDAGJSON, _ := json.Marshal(serverlessDAG)
	err = create(serverlessDAGJSON, service_const.ServerlessDAGPrefix, serverlessDAG.Metadata.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, serverlessDAG)
}

func updateServerlessDAG(c *gin.Context) {
	var serverlessDAG Serverless.DAG
	err := c.BindJSON(&serverlessDAG)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	serverlessDAGJSON, _ := json.Marshal(serverlessDAG)
	err = put(serverlessDAGJSON, service_const.ServerlessDAGPrefix, serverlessDAG.Metadata.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, serverlessDAG)
}

func deleteServerlessDAG(c *gin.Context) {
	serverlessDAGName := c.Param("name")
	if serverlessDAGName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	err := delete(service_const.ServerlessDAGPrefix, serverlessDAGName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func listServerlessDAG(c *gin.Context) {
	ifWatch := c.DefaultQuery("watch", "false")

	// short HTTP connection, only list the serverlessDAGs
	if ifWatch == "false" {
		serverlessDAGList := make([]Serverless.DAG, 0)
		serverlessDAGJSONList, err := list(service_const.ServerlessDAGPrefix)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		for _, serverlessDAGJSON := range serverlessDAGJSONList {
			var serverlessDAG Serverless.DAG
			err = json.Unmarshal([]byte(serverlessDAGJSON), &serverlessDAG)
			if err != nil {
				continue
			}
			serverlessDAGList = append(serverlessDAGList, serverlessDAG)
		}

		c.JSON(http.StatusOK, serverlessDAGList)
		return
	}

	// long HTTP connection, create a watcher and start watching
	watchList(c, service_const.ServerlessDAGPrefix)
}

func getServerlessDAG(c *gin.Context) {
	serverlessDAGName := c.Param("name")
	if serverlessDAGName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	serverlessDAGJSON, err := get(service_const.ServerlessDAGPrefix, serverlessDAGName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if serverlessDAGJSON != "" {
		var serverlessDAG Serverless.DAG
		err = json.Unmarshal([]byte(serverlessDAGJSON), &serverlessDAG)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, serverlessDAG)
		return
	}

	c.JSON(http.StatusOK, nil)
}
