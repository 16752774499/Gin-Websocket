package model

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	UserName  string `gorm:"unique"`
	Password  string `gorm:"not null"`
	Avatar    sql.NullString
	DeletedAt gorm.DeletedAt `gorm:"index"` // 带索引的软删除时间字段，用于软删除标记

	//设置关联外键
	// 定义一对多关系，一个用户有多个好友关系（反向关联）
	FriendRelationships []UserFriend `gorm:"foreignKey:UserID"`
}

// 加密难度
const PassWordCost = 12

// SetPassword 设置用户的密码
//
// 参数:
//
//	password: 用户的新密码，类型为字符串
//
// 返回值:
//
//	如果密码设置成功，则返回 nil；如果设置失败，则返回错误信息
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(password string) bool {
	// 使用 bcrypt 库比较用户存储的密码哈希值和输入的密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	// 返回比较结果，如果比较成功则返回 true，否则返回 false
	return err == nil
}
