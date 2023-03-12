package server

import (
	"encoding/json"
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/data/Serverless"
	"example/Minik8s/pkg/serverless/data"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

var DAGMap map[string]*Serverless.DAG

func StartServerlessDAG(serverlessDAG Serverless.DAG) {
	DAGName := serverlessDAG.Metadata.Name
	_, ok := DAGMap[DAGName]
	if ok {
		log.Println("DAG " + DAGName + " already exists!")
		return
	}

	serverlessDAG.Status.Phase = "running"
	for _, step := range serverlessDAG.Spec.Steps {
		if step.Type == "task" {
			StartServerlessService(step.Task.Function)
		}
	}

	DAGMap[DAGName] = &serverlessDAG

	apiclient.Request(data.GetCredential(), "/api/v1/serverlessDAG/", serverlessDAG, "PUT")
}

func handleDAGRequest(c *gin.Context) {
	DAGName := c.Param("name")
	if DAGName == "" {
		log.Println("empty DAG name")
		c.JSON(http.StatusBadRequest, "empty DAG name")
		return
	}

	serverlessDAG, ok := DAGMap[DAGName]
	if !ok {
		log.Println("no such DAG")
		c.JSON(http.StatusBadRequest, "no such DAG")
		return
	}

	query := make(map[string]string)
	for _, key := range serverlessDAG.Spec.Input {
		val := c.DefaultQuery(key, "")
		query[key] = val
	}

	var resBody []byte
	response := new(http.Response)
	var err error
	for i := 0; i < len(serverlessDAG.Spec.Steps); {
		step := serverlessDAG.Spec.Steps[i]
		log.Println("step " + step.Name)
		if step.Type == "task" {
			serverlessService, ok := serviceMap[step.Task.Function.Metadata.Name]
			if !ok {
				log.Println("no such service")
				c.JSON(http.StatusBadRequest, "no such service")
				return
			}
			response, err = forwardRequest(serverlessService, query)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			resBody, _ = ioutil.ReadAll(response.Body)
			resMap := make(map[string]string)
			err = json.Unmarshal(resBody, &resMap)
			if err != nil {
				log.Println(err)
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			query = make(map[string]string)
			for _, key := range step.Task.Function.Spec.Output {
				val := resMap[key]
				query[key] = val
			}
			i = i + 1
		} else if step.Type == "choice" {
			key := step.Choice.Key
			val, ok := query[key]
			if !ok {
				log.Println("no such key")
				c.JSON(http.StatusBadRequest, "so such key")
				return
			}
			ifMatch := false
			for targetVal, targetFunc := range step.Choice.Jump {
				if val == targetVal {
					ifMatch = true
					var j int
					for j = 0; j < len(serverlessDAG.Spec.Steps); j++ {
						if serverlessDAG.Spec.Steps[j].Name == targetFunc {
							break
						}
					}
					if j >= len(serverlessDAG.Spec.Steps) {
						log.Println("error jump target")
						c.JSON(http.StatusBadRequest, "error jump target")
						return
					}
					i = j
					break
				}
			}
			if ifMatch == false {
				i = i + 1
			}
		} else {
			log.Println("unsupported type")
			c.JSON(http.StatusBadRequest, "unsupported type")
			return
		}
	}

	c.Data(response.StatusCode, response.Header.Get("content-type"), resBody)
}
