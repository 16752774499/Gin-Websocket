package router

import (
	"Gin-WebSocket/api"
	"github.com/gin-gonic/gin"
)

func FileRouter(router *gin.Engine) {
	fileRouter := router.Group("/HandleFile")
	{
		fileRouter.POST("/upload", api.Upload)
	}
}
