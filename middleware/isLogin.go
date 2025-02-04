package middleware

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func IsLogin(ctx *gin.Context) {
	session := sessions.Default(ctx)
	ret := session.Get("userInfo")
	if ret == nil {
		//走登录流程
		ctx.Next()
	} else {
		fmt.Println(ret)
		return
	}
}
