package api

import (
	"Gin-WebSocket/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func UserRegister(ctx *gin.Context) {
	var userRegisterService service.UserRegisterService
	if err := ctx.ShouldBind(&userRegisterService); err != nil {
		logrus.Info("UserRegister err: ", err)
		ctx.JSON(400, ErrorResponse(err))

	} else {
		res := userRegisterService.Register()
		ctx.JSON(200, res)

	}

}
