package main

import (
	"Gin-WebSocket/cache"
	"Gin-WebSocket/conf"
	"Gin-WebSocket/middleware"
	"Gin-WebSocket/router"
	"Gin-WebSocket/service/wsHeartbeat"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func main() {
	conf.Init()
	//go service.StartChat()
	go wsHeartbeat.StartHeartbeats()
	r := gin.Default()
	r.Static("/static", "./statics")
	r.Use(gin.Recovery(), gin.Logger(), middleware.CORSMiddleware())

	r.Use(sessions.Sessions("Mysession", cache.NewSessionStore()))
	router.UserRouter(r)
	router.WsRouter(r)
	err := r.Run(conf.HttpPort)
	if err != nil {
		panic(err)
	}

}
