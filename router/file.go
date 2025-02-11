package router

import (
	"Gin-WebSocket/api"

	"github.com/gin-gonic/gin"
)

func FileRouter(router *gin.Engine) {
	fileRouter := router.Group("/file")
	{
		fileRouter.POST("/upload", api.Upload)
	}
}
