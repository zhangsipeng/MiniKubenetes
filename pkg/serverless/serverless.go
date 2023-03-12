package serverless

import (
	"example/Minik8s/pkg/apiclient"
	"example/Minik8s/pkg/serverless/data"
	"example/Minik8s/pkg/serverless/server"
	"example/Minik8s/pkg/serverless/watchApiServer"
	"github.com/gin-gonic/gin"
)

func StartService() {
	info := apiclient.GetInitInfo()
	runtimeConfig, err := apiclient.InitRuntimeConfigOrError(info, nil)
	if err != nil {
		panic(err)
	}
	router := gin.Default()
	server.InitRouter(router)
	server.InitServer()
	data.InitCredential(runtimeConfig)
	go watchApiServer.WatchServerlessService()
	go watchApiServer.WatchServerlessDAG()
	err = router.Run(":50000")
	if err != nil {
		panic(err)
	}
}
