package service

import "github.com/gin-gonic/gin"

func RegisterTLSService(router *gin.Engine) {
	// pod
	router.POST("/api/v1/pods/", createPod)
	router.PUT("/api/v1/pods/", updatePod)
	router.GET("/api/v1/pods/", listPod)
	router.GET("/api/v1/pods/:name/", getPod)
	router.DELETE("/api/v1/pods/:name/", deletePod)
	router.POST("/api/v1/podInNode/:name/pods/", createPodInNode)
	router.DELETE("/api/v1/podInNode/:nodeName/pods/:podName/", deletePodInNode)
	router.GET("/api/v1/podInNode/:name/pods/", listPodInNode)

	// node
	router.POST("/api/v1/nodes/", createNode)
	router.PUT("/api/v1/nodes/", updateNode)
	router.GET("/api/v1/nodes/", listNode)
	router.GET("/api/v1/nodes/:name/", getNode)
	router.DELETE("/api/v1/nodes/:name/", deleteNode)

	// service
	router.POST("/api/v1/service/", createService)
	router.PUT("/api/v1/service/", updateService)
	router.GET("/api/v1/service/", listService)
	router.GET("/api/v1/service/:name/", getService)
	router.DELETE("/api/v1/service/:name/", deleteService)

	// endpoints
	router.POST("/api/v1/endpoints/", createEndpoints)
	router.PUT("/api/v1/endpoints/", updateEndpoints)
	router.GET("/api/v1/endpoints/", listEndpoints)
	router.GET("/api/v1/endpoints/:name/", getEndpoints)
	router.DELETE("/api/v1/endpoints/:name/", deleteEndpoints)

	// horizontalPodAutoscaler
	router.POST("/api/v1/horizontalPodAutoscaler", createHorizontalPodAutoscaler)
	router.PUT("/api/v1/horizontalPodAutoscaler", updateHorizontalPodAutoscaler)
	router.GET("/api/v1/horizontalPodAutoscaler", listHorizontalPodAutoscaler)
	router.GET("/api/v1/horizontalPodAutoscaler/:name", getHorizontalPodAutoscaler)
	router.DELETE("/api/v1/horizontalPodAutoscaler", deleteHorizontalPodAutoscaler)

	// job
	router.POST("/api/v1/gpujobs/", createJob)
	router.PUT("/api/v1/gpujobs/", updateJob)
	router.GET("/api/v1/gpujobs/", listJob)
	router.GET("/api/v1/gpujobs/:name/", getJob)
	router.DELETE("/api/v1/gpujobs/:name/", deleteJob)
	router.POST("/api/v1/jobInNode/:name/jobs/", createJobInNode)
	router.DELETE("/api/v1/jobInNode/:nodeName/jobs/:podName/", deleteJobInNode)
	router.GET("/api/v1/jobInNode/:name/jobs/", listJobInNode)

	// replicaSet
	router.POST("/api/v1/replicaSet/", createReplicaSet)
	router.PUT("/api/v1/replicaSet/", updateReplicaSet)
	router.GET("/api/v1/replicaSet/", listReplicaSet)
	router.GET("/api/v1/replicaSet/:name/", getReplicaSet)
	router.DELETE("/api/v1/replicaSet/:name/", deleteReplicaSet)

	// deployment
	router.POST("/api/v1/deployment/", createDeployment)
	router.PUT("/api/v1/deployment/", updateDeployment)
	router.GET("/api/v1/deployment/", listDeployment)
	router.GET("/api/v1/deployment/:name/", getDeployment)
	router.DELETE("/api/v1/deployment/:name/", deleteDeployment)

	// serverless
	router.POST("/api/v1/serverless/", createServerlessService)
	router.PUT("/api/v1/serverless/", updateServerlessService)
	router.GET("/api/v1/serverless/", listServerlessService)
	router.GET("/api/v1/serverless/:name/", getServerlessService)
	router.DELETE("/api/v1/serverless/:name/", deleteServerlessService)

	// serverlessDAG
	router.POST("/api/v1/serverlessDAG/", createServerlessDAG)
	router.PUT("/api/v1/serverlessDAG/", updateServerlessDAG)
	router.GET("/api/v1/serverlessDAG/", listServerlessDAG)
	router.GET("/api/v1/serverlessDAG/:name/", getServerlessDAG)
	router.DELETE("/api/v1/serverlessDAG/:name/", deleteServerlessDAG)
}

func RegisterNonTLSService(router *gin.Engine) {
	// register token service
	router.POST("/api/v1/token/", checkToken)
}
