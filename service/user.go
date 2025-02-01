package service

import (
	"Gin-WebSocket/model"
	"Gin-WebSocket/serializer"
)

type UserRegisterService struct {
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
}

func (service *UserRegisterService) Register() serializer.Response {
	var user model.User
	var count int64 = 0
	/*
		.Model(&model.User{}): 这个方法告诉GORM我们想要操作的模型是model.User。
		&model.User{}是一个指向User结构体实例的指针，它表示了我们想要查询或操作的模型类型。
		在这个上下文中，它并不表示具体的某个用户，而是告诉GORM我们接下来要进行的操作是针对User表的。
	*/
	model.DB.Model(&model.User{}).Where("user_name = ?", service.UserName).First(&user).Count(&count)
	if count != 0 {

		return serializer.Response{
			Status: 400,
			Msg:    "用户名已经存在！",
		}
	}
	user = model.User{
		UserName: service.UserName,
	}
	if err := user.SetPassword(service.Password); err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "密码加密出错！",
		}
	}
	model.DB.Create(&user)
	return serializer.Response{
		Status: 200,
		Msg:    "创建成功！",
	}
}
