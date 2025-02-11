package api

import (
	"Gin-WebSocket/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	MaxFileSize = 100 << 20 // 100MB
)

// 文件上传
func Upload(ctx *gin.Context) {

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": "No HandleFile uploaded"})
		logrus.Error(err)
		return
	}
	defer file.Close()
	// 检查文件大小
	if header.Size > MaxFileSize {
		ctx.JSON(http.StatusBadRequest, service.ErrorResponse{
			Status:  400,
			Message: "File size exceeds limit (100MB)",
		})
		//ctx.JSON(http.StatusBadRequest, gin.H{
		//	"status": 400,
		//	"msg":    "File size exceeds limit (100MB)",
		//})
		return
	}
	res, err := service.UploadFile(file, header)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, service.ErrorResponse{
			Status:  400,
			Message: err.Error(),
			Error:   "handleFile-UploadFile error!",
		})
		return
	}
	ctx.JSON(http.StatusOK, res)

}
