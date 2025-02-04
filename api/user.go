package api

import (
	"Gin-WebSocket/model"
	"Gin-WebSocket/serializer"
	"Gin-WebSocket/service"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strconv"
)

func UserRegister(ctx *gin.Context) {
	var userRegisterService service.UserService
	if err := ctx.ShouldBind(&userRegisterService); err != nil {
		logrus.Info("UserRegister err: ", err)
		ctx.JSON(400, ErrorResponse(err))

	} else {
		res := userRegisterService.Register()
		ctx.JSON(200, res)

	}

}

func UserLogin(ctx *gin.Context) {
	var userLoginService service.UserService
	if err := ctx.ShouldBind(&userLoginService); err != nil {
		logrus.Info("Login err: ", err)
		ctx.JSON(400, ErrorResponse(err))
	} else {
		res := userLoginService.Login(ctx)
		ctx.JSON(200, res)
	}
}

func User(ctx *gin.Context) {

	if ctx.Request.Method == "GET" {

	} else if ctx.Request.Method == "POST" {
		var userAvatarService service.UserService
		if err := ctx.ShouldBind(&userAvatarService); err != nil {
			logrus.Info("USER POST err: ", err)
			ctx.JSON(400, ErrorResponse(err))
		} else {
			res := userAvatarService.POST()
			ctx.JSON(200, res)
		}

	} else if ctx.Request.Method == "PUT" {

	} else if ctx.Request.Method == "DELETE" {

	}

}
func CheckSession(ctx *gin.Context) {

	session := sessions.Default(ctx)
	ret := session.Get("userInfo")
	if ret == nil {
		//没有session,去登录
		//logrus.Info("Login err: ", err)

		ctx.JSON(400, serializer.Response{
			Status: 400,
			Msg:    "无效Session！",
		})
	} else {
		ok, res := service.CheckSession(ret)
		if ok {
			ctx.JSON(200, res)
		} else {
			ctx.JSON(400, res)
		}
	}
}

func SearchUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userName := ctx.PostForm("user_name")
	//如果搜索用户名为空或是自己

	if userName == "" {
		ctx.JSON(400, serializer.Response{
			Status: 400,
			Msg:    "好友名称不合法！",
		})
	} else if userName == session.Get("userInfo") {
		ctx.JSON(200, serializer.Response{
			Status: 400,
			Msg:    "不能添加自己为好友！",
		})
	} else {
		res := service.SearchUser(userName)
		ctx.JSON(200, res)

	}
}

func AddFriend(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userName := session.Get("userInfo")
	//获取用户id
	var user model.User
	model.DB.Where("user_name = ?", userName).First(&user)

	friendIdString := ctx.PostForm("friend_id")
	friendId, err := strconv.Atoi(friendIdString)
	if err != nil {
		logrus.Info("无法将 friend_id 转换为 int 类型: %v", err)
		ctx.JSON(400, serializer.Response{
			Status: 400,
			Msg:    fmt.Sprintf("无法将 friend_id 转换为 int 类型: %v", err),
		})
		return
	}
	logrus.Info("friendId: ", friendId)
	res := service.AddFriend(int(user.ID), friendId)
	ctx.JSON(200, res)
}

func FriendRequests(ctx *gin.Context) {
	var user model.User
	session := sessions.Default(ctx)
	userName := session.Get("userInfo")
	model.DB.Where("user_name = ?", userName).First(&user)
	res := service.FriendRequests(user.ID)
	ctx.JSON(200, res)
}

func HandleRequest(ctx *gin.Context) {
	requestId := ctx.PostForm("request_id")
	action := ctx.PostForm("action")
	if ok, res := service.HandleRequest(requestId, action); ok {
		ctx.JSON(200, res)
	} else {
		ctx.JSON(400, res)
	}

}

func Friend(context *gin.Context) {
	session := sessions.Default(context)
	userName := session.Get("userInfo")
	var user model.User
	model.DB.Where("user_name = ?", userName).First(&user)
	res := service.Friends(user.ID)
	context.JSON(200, res)

}
