package router

import (
	"Gin-WebSocket/api"
	"github.com/gin-gonic/gin"
)

//func UserRouter() *gin.RouterGroup {
//	r := gin.Default()
//	/*  Recovery 中间件
//	作用：Recovery 中间件用于从任何 panics 中恢复，
//	并确保 Gin 应用程序不会因为未处理的 panics 而崩溃。
//	当 Gin 应用程序在处理某个请求时发生 panic，
//	Recovery 中间件会捕获这个 panic，记录相关错误信息（通常记录到日志中），
//	并且返回一个合适的 HTTP 响应给客户端，例如返回一个 500 Internal Server Error 的 HTTP 状态码。
//	*/
//	/*
//		在 Go 语言中，panics（复数形式，单数是 panic）是一种运行时错误处理机制，
//		用于表示程序遇到了严重的、无法正常恢复的问题。当 panic 发生时，
//		程序会立即停止当前函数的执行，并开始展开堆栈（unwind the stack）。
//		即从调用栈中从当前函数依次向上移除每层调用，释放	每层调用中分配的局部变量，
//		直到找到对应的 recover 函数或者程序终止。
//	*/
//	r.Use(gin.Recovery(), gin.Logger())
//	v1 := r.Group("/")
//	{
//
//		v1.POST("user/register", api.UserRegister)
//		v1.POST("user/login", api.UserLogin)
//	}
//	return r.Group()
//}

func UserRouter(router *gin.Engine) {
	userRouter := router.Group("/user")
	{
		userRouter.POST("/register", api.UserRegister)
		userRouter.POST("/login", api.UserLogin)
		userRouter.GET("/check-session", api.CheckSession)
		userRouter.POST("/search", api.SearchUser)    //搜索用户
		userRouter.POST("/add-friend", api.AddFriend) //添加好友
		userRouter.GET("/friend-requests", api.FriendRequests)
		userRouter.POST("handle-request", api.HandleRequest)
		userRouter.GET("/friends", api.Friend)

		userRouter.POST("", api.User)             // 用户信息
		userRouter.DELETE("/:userName", api.User) //注销用户
		userRouter.PUT("", api.User)              //修改用户信息

	}
}
