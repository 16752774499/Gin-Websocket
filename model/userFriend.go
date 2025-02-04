package model

import (
	"gorm.io/gorm"
)

type UserFriend struct {
	// ID 字段作为好友关系表的主键，用于唯一标识每一条好友关系记录
	ID uint `gorm:"primaryKey"`
	// UserID 字段是外键，关联到 User 表的主键 ID
	// 它表示发起好友关系的用户 ID
	UserID uint
	// FriendID 字段是外键，关联到 User 表的主键 ID
	// 它表示被添加为好友的用户 ID
	FriendID uint
	// IsAccepted 字段用于标识好友申请是否被接受
	// 值为 true 表示已接受，false 表示未接受
	IsAccepted bool
	// RequestTime 字段以时间戳（int64 类型）记录好友申请发送的时间
	RequestTime int64
	// AcceptedTime 字段以指针形式存储好友申请被接受的时间戳
	// 指针类型意味着该字段在数据库中可以为 NULL，表示申请尚未被接受
	AcceptedTime *int64

	DeletedAt gorm.DeletedAt `gorm:"index"` // 带索引的软删除时间字段，用于软删除标记
	// User 字段用于定义与 User 结构体的关联关系
	// gorm 标签中的 references:ID 表示引用 User 结构体中的 ID 字段作为外键关联依据
	User User `gorm:"references:ID"`
}

//func (user *User) (password string) bool {
//	// 使用 bcrypt 库比较用户存储的密码哈希值和输入的密码
//	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
//	// 返回比较结果，如果比较成功则返回 true，否则返回 false
//	return err == nil
//}
