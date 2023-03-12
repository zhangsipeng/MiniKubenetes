package service

import (
	"encoding/json"
	"example/Minik8s/pkg/data/WorkloadResources"
	"example/Minik8s/pkg/kubeapiserver/service/const"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func createDeployment(c *gin.Context) {
	var deployment WorkloadResources.Deployment
	err := c.BindJSON(&deployment)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	deploymentJSON, _ := json.Marshal(deployment)
	err = create(deploymentJSON, service_const.DeploymentPrefix, deployment.Metadata.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, deployment)
}

func updateDeployment(c *gin.Context) {
	var deployment WorkloadResources.Deployment
	err := c.BindJSON(&deployment)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	deploymentJSON, _ := json.Marshal(deployment)
	err = put(deploymentJSON, service_const.DeploymentPrefix, deployment.Metadata.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, deployment)
}

func deleteDeployment(c *gin.Context) {
	deploymentName := c.Param("name")
	if deploymentName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	err := delete(service_const.DeploymentPrefix, deploymentName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

func listDeployment(c *gin.Context) {
	ifWatch := c.DefaultQuery("watch", "false")

	// short HTTP connection, only list the deployments
	if ifWatch == "false" {
		deploymentList := make([]WorkloadResources.Deployment, 0)
		deploymentJSONList, err := list(service_const.DeploymentPrefix)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		for _, deploymentJSON := range deploymentJSONList {
			var deployment WorkloadResources.Deployment
			err = json.Unmarshal([]byte(deploymentJSON), &deployment)
			if err != nil {
				continue
			}
			deploymentList = append(deploymentList, deployment)
		}

		c.JSON(http.StatusOK, deploymentList)
		return
	}

	// long HTTP connection, create a watcher and start watching
	watchList(c, service_const.DeploymentPrefix)
}

func getDeployment(c *gin.Context) {
	deploymentName := c.Param("name")
	if deploymentName == "" {
		log.Println(service_const.EmptyNameError)
		c.JSON(http.StatusBadRequest, service_const.EmptyNameError)
		return
	}

	deploymentJSON, err := get(service_const.DeploymentPrefix, deploymentName)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if deploymentJSON != "" {
		var deployment WorkloadResources.Deployment
		err = json.Unmarshal([]byte(deploymentJSON), &deployment)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		c.JSON(http.StatusOK, deployment)
		return
	}

	c.JSON(http.StatusOK, nil)
}
