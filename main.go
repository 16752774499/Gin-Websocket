package main

import (
	"Gin-WebSocket/conf"
	"Gin-WebSocket/router"
	"Gin-WebSocket/service"
)

func main() {
	conf.Init()
	go service.Manager.Start()
	r := router.NewRouter()
	err := r.Run(conf.HttpPort)
	if err != nil {
		panic(err)
	}

}
