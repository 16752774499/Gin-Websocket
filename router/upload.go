package router

import (
	"Gin-WebSocket/api"
	"github.com/gin-gonic/gin"
)

func UploadRouter(router *gin.Engine) {
	uploadRouter := router.Group("/upload")
	{
		uploadRouter.POST("/upload", api.Upload)
	}
}
