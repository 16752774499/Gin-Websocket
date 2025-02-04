package tools

import (
	"Gin-WebSocket/serializer"
	"time"
)

func NewTimeStamp() *int64 {
	now := time.Now().Unix()
	return &now
}

func ServiceError() serializer.Response {
	return serializer.Response{
		Status: 500,
		Msg:    "服务器异常",
	}
}
