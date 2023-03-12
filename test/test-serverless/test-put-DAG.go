package main

import (
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/data/ObjectMeta"
	"example/Minik8s/pkg/data/Serverless"
	"log"
)

func main() {
	serverlessDAG := Serverless.DAG{
		ApiVersion: "v1",
		Kind:       "serverlessDAG",
		Metadata:   ObjectMeta.ObjectMeta{Name: "test"},
		Spec: Serverless.DAGSpec{
			Steps: []Serverless.Step{
				{
					Name: "add",
					Type: "task",
					Task: Serverless.StepTask{Function: Serverless.Service{
						ApiVersion: "v1",
						Kind:       "serverlessService",
						Metadata:   ObjectMeta.ObjectMeta{Name: "add"},
						Spec: Serverless.ServiceSpec{
							GitUrl:      "/root/serverless/add",
							MaxReplicas: 5,
							Input:       []string{"A", "B"},
							Output:      []string{"A"},
						},
						Status: Serverless.ServiceStatus{Phase: "init"},
					}},
				},
				{
					Name: "compareZero",
					Type: "choice",
					Choice: Serverless.StepChoice{
						Key: "A",
						Jump: map[string]string{
							"1":  "reverse",
							"2":  "reverse",
							"3":  "reverse",
							"4":  "reverse",
							"5":  "reverse",
							"6":  "reverse",
							"7":  "reverse",
							"8":  "reverse",
							"9":  "reverse",
							"10": "reverse",
						},
					},
				},
				{
					Name: "square",
					Type: "task",
					Task: Serverless.StepTask{Function: Serverless.Service{
						ApiVersion: "v1",
						Kind:       "serverlessService",
						Metadata:   ObjectMeta.ObjectMeta{Name: "square"},
						Spec: Serverless.ServiceSpec{
							GitUrl:      "/root/serverless/square",
							MaxReplicas: 5,
							Input:       []string{"A"},
							Output:      []string{"A"},
						},
						Status: Serverless.ServiceStatus{Phase: "init"},
					}},
				},
				{
					Name: "reverse",
					Type: "task",
					Task: Serverless.StepTask{Function: Serverless.Service{
						ApiVersion: "v1",
						Kind:       "serverlessService",
						Metadata:   ObjectMeta.ObjectMeta{Name: "reverse"},
						Spec: Serverless.ServiceSpec{
							GitUrl:      "/root/serverless/reverse",
							MaxReplicas: 5,
							Input:       []string{"A"},
							Output:      []string{"A"},
						},
						Status: Serverless.ServiceStatus{Phase: "init"},
					}},
				},
			},
			Input: []string{"A", "B"},
		},
		Status: Serverless.DAGStatus{Phase: "init"},
	}

	info := apiclient.GetInitInfo()
	runtimeConfig, _ := apiclient.InitRuntimeConfigOrError(info, nil)

	response := apiclient.Request(runtimeConfig, "/api/v1/serverlessDAG/", serverlessDAG, "POST")

	log.Println(string(response))
}
