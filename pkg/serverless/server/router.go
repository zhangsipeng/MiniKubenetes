package server

import "github.com/gin-gonic/gin"

var router *gin.Engine

func InitRouter(r *gin.Engine) {
	router = r
}

func AddService(method, relativePath string, handler func(ctx *gin.Context)) {
	router.Handle(method, relativePath, handler)
}
