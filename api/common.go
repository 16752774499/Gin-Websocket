package api

import (
	"Gin-WebSocket/serializer"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
)

func ErrorResponse(err error) serializer.Response {
	/*
		validator库是 Go 语言中用于数据验证的一个强大工具，主要用于对结构体字段、基本类型等数据进行各种规则的验证。
		在使用 Go 编写 Web 服务（如 Gin 框架的应用）时，validator/v10 库常用于验证用户输入的数据，
		确保接收到的数据符合预期的格式和限制。这有助于防止因不合理的输入导致程序出现错误或者安全风险。
	*/
	// 检查错误是否为 validator.ValidationErrors 类型
	if _, ok := err.(validator.ValidationErrors); ok {
		// 返回状态码为 400 的响应，表示参数错误
		return serializer.Response{
			Status: 400,
			Msg:    "参数错误！",
			Error:  fmt.Sprintf("%v", err),
		}
	}
	// 检查错误是否为 json.UnmarshalTypeError 类型
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		// 返回状态码为 400 的响应，表示 JSON 类型不匹配
		return serializer.Response{
			Status: 400,
			Msg:    "JSON类型不匹配！",
			Error:  fmt.Sprintf("%v", err),
		}
	}
	//返回状态码为 500 的响应，表示其他类型的错误
	return serializer.Response{
		Status: 400,
		Msg:    "参数错误！",
		Error:  fmt.Sprintf("%v", err),
	}

}
