package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Upload(ctx *gin.Context) {
	fmt.Println(ctx.Request)
}
