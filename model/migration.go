package model

import "fmt"

func migration() {

	//自动迁移表结构
	err = DB.AutoMigrate(&User{})
	if err != nil {
		fmt.Printf("Failed to auto migrate: %v\n", err)
	}
}
