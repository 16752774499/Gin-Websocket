package model

import "fmt"

func Migration() {

	//自动迁移表结构
	err = DB.AutoMigrate(&User{}, UserFriend{})
	if err != nil {
		fmt.Printf("Failed to auto migrate: %v\n", err)
	}
}
