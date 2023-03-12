package service

import (
	"encoding/json"
	"example/Minik8s/pkg/data/ClusterResources"
	service_const "example/Minik8s/pkg/kubeapiserver/service/const"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func createNode(c *gin.Context) {
	var node ClusterResources.Node
	err := c.BindJSON(&node)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	nodeJSON, _ := json.Marshal(node)
	err = create(nodeJSON, service_const.NodePrefix, node.Metadata.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, node)
}

func updateNode(c *gin.Context) {
	var node ClusterResources.Node
	err := c.BindJSON(&node)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	nodeJSON, _ := json.Marshal(node)
	err = put(nodeJSON, service_const.NodePrefix, node.Metadata.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, node)
}

func deleteNode(c *gin.Context) {
	nodeName := c.Param("name")
	if nodeName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	err := delete(service_const.NodePrefix, nodeName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func listNode(c *gin.Context) {
	ifWatch := c.DefaultQuery("watch", "false")

	// short HTTP connection, only list the nodes
	if ifWatch == "false" {
		nodeList := make([]ClusterResources.Node, 0)
		nodeJSONList, err := list(service_const.NodePrefix)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		for _, nodeJSON := range nodeJSONList {
			var node ClusterResources.Node
			err = json.Unmarshal([]byte(nodeJSON), &node)
			if err != nil {
				continue
			}
			nodeList = append(nodeList, node)
		}

		c.JSON(http.StatusOK, nodeList)
		return
	}

	// long HTTP connection, create a watcher and start watching
	watchList(c, service_const.NodePrefix)
}

func getNode(c *gin.Context) {
	nodeName := c.Param("name")
	if nodeName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	nodeJSON, err := get(service_const.NodePrefix, nodeName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if nodeJSON != "" {
		var node ClusterResources.Node
		err = json.Unmarshal([]byte(nodeJSON), &node)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, node)
		return
	}

	c.JSON(http.StatusOK, nil)
}
