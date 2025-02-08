package handleUser

import (
	"Gin-WebSocket/model"
	tools "Gin-WebSocket/public"
	"Gin-WebSocket/serializer"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type UserService struct {
	Id       int    `json:"id"`
	UserName string `json:"user_name" form:"user_name"`
	Password string `json:"password" form:"password"`
}

func (service *UserService) Register() serializer.Response {
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

func (service *UserService) Login(ctx *gin.Context) serializer.Response {
	session := sessions.Default(ctx)
	userName := service.UserName
	var user model.User
	model.DB.Model(&model.User{}).Where("user_name = ?", userName).First(&user)
	if user.UserName != service.UserName {
		return serializer.Response{
			Status: 400,
			Msg:    "账号不存在，请前往注册！",
		}
	} else {
		//检查密码
		if user.CheckPassword(service.Password) {
			data := map[string]interface{}{
				"id":       user.ID,
				"userName": user.UserName,
			}
			//账号密码正确，设置Session
			session.Set("userInfo", user.UserName)
			err := session.Save()
			if err != nil {
				return serializer.Response{
					Status: 500,
					Msg:    "Seesion 加载失败！",
					Data:   err.Error(),
				}
			}
			return serializer.Response{
				Status: 200,
				Msg:    "登录成功",
				Data:   data,
			}
		} else {
			return serializer.Response{
				Status: 400,
				Msg:    "账号或密码错误！",
			}
		}
	}

}

func (service *UserService) POST() serializer.Response {
	userID := service.Id
	var user model.User
	model.DB.Model(&model.User{}).Where("id = ?", userID).First(&user)

	return serializer.Response{
		Status: 200,
		Msg:    "请求成功！",
		Data:   user.Avatar.String,
	}
}

func CheckSession(userInfo interface{}) (bool, serializer.Response) {
	var user model.User
	model.DB.Model(&model.User{}).Where("user_name = ?", userInfo).First(&user)
	if user.ID == 0 {
		return false, serializer.Response{
			Status: 400,
			Msg:    "无效Session",
		}
	} else {
		data := map[string]interface{}{
			"id":       user.ID,
			"userName": user.UserName,
			"avatar":   user.Avatar.String,
		}
		return true, serializer.Response{
			Status: 200,
			Data:   data,
			Msg:    "session有效",
		}
	}

}

func SearchUser(name string) serializer.Response {
	//查询库中有没有该用户
	var user model.User
	model.DB.Model(&model.User{}).Where("user_name = ?", name).First(&user)
	if user.ID == 0 {
		//没有该用户
		return serializer.Response{
			Status: 400,
			Msg:    "该用户不存在！",
			Data:   nil,
		}
	} else {
		data := map[string]interface{}{
			"userName": user.UserName,
			"id":       user.ID,
			"avatar":   user.Avatar.String,
		}
		return serializer.Response{
			Status: 200,
			Data:   data,
			Msg:    "查找成功！",
		}
	}

}

// 检查反向关系
func AddFriend(userId int, friendId int) serializer.Response {

	//检查是不是已经是好友了
	var userFriend model.UserFriend
	model.DB.Where("user_id = ? AND  friend_id = ?", uint(userId), uint(friendId)).First(&userFriend)
	if userFriend.ID == 0 {
		//还不是好友
		logrus.Info("加好友！！")
		newFriendRequest := model.UserFriend{
			UserID:      uint(userId),
			FriendID:    uint(friendId),
			IsAccepted:  false,
			RequestTime: time.Now().Unix(), // 实际应用中应替换为当前时间戳
		}
		result := model.DB.Create(&newFriendRequest)
		if result.Error != nil {
			return serializer.Response{
				Status: 500,
				Msg:    fmt.Sprintf("请求失败:%s", result.Error.Error()),
			}
		} else {
			logrus.Info("请求成功！:")
			return serializer.Response{
				Status: 200,
				Msg:    "好友申请已提交！",
			}
		}
	} else {
		//判断是否通过
		if userFriend.IsAccepted {
			return serializer.Response{
				Status: 400,
				Msg:    "你们已经是好友了！",
			}
		} else {
			return serializer.Response{
				Status: 400,
				Msg:    "请勿重复提交申请！",
			}
		}

	}

}

func FriendRequests(userId uint) serializer.Response {
	type From struct {
		Id       string `json:"id"`
		UserName string `json:"userName"`
		Avatar   string `json:"avatar"`
	}
	type FriendRequest struct {
		ID        int       `json:"id"`
		From      From      `json:"from"`
		CreatedAt time.Time `json:"createdAt"`
	}

	Data := []FriendRequest{}

	tx := model.DB.Begin()
	var friendRequests []model.UserFriend
	err := tx.Where("friend_id =? AND is_accepted =?", userId, 0).Find(&friendRequests).Error
	//无申请
	if len(friendRequests) == 0 {
		return serializer.Response{
			Status: 200,
			Data:   nil,
			Msg:    "当前无好友申请！",
		}
	}
	//双方已经有一条
	if err != nil {
		tx.Rollback()
		return serializer.Response{
			Status: 500,
			Msg:    fmt.Sprintf("获取好友申请出错:%s", err.Error()),
		}
	} else {
		//
		tx.SavePoint("xiao") //创建回滚标记
		for _, friendRequest := range friendRequests {
			data := FriendRequest{}
			data.ID = int(friendRequest.ID)
			data.CreatedAt = time.Unix(friendRequest.RequestTime, 0)
			from := From{}
			err := model.DB.Model(&model.User{}).Select("id", "user_name", "avatar").Where("id = ?", friendRequest.UserID).First(&from).Error
			if err != nil {
				tx.RollbackTo("xiao")
				return serializer.Response{
					Status: 500,
					Msg:    fmt.Sprintf("获取好友信息出错:%s", err.Error()),
				}
			} else {
				data.From = from
			}
			Data = append(Data, data)
		}
		return serializer.Response{
			Status: 200,
			Msg:    "获取好友请求列表成功",
			Data:   Data,
		}
	}

}

func HandleRequest(requestId string, action string) (bool, serializer.Response) {
	if action == "reject" {
		//不同意 ，删除此条记录
		tx := model.DB.Begin()
		//软删除该条数据据
		//result := models.DB.Where("id = ?", id).Delete(&models.User{})

		result := tx.Where("id = ?", requestId).Delete(&model.UserFriend{})
		if result.Error != nil {
			tx.Rollback()
			return false, serializer.Response{
				Status: 500,
				Msg:    "服务器异常！",
			}
		} else {
			tx.Commit()
			return true, serializer.Response{
				Status: 200,
				Msg:    "已拒绝好友申请!",
			}
		}
	} else if action == "accept" {
		//同意
		tx := model.DB.Begin()
		var userFriend model.UserFriend
		//获取该条id记录所对应的两位用户
		result := tx.Where("id = ?", requestId).First(&userFriend)
		if result.Error != nil {
			tx.Rollback()
			return false, serializer.Response{
				Status: 500,
				Msg:    "服务器异常！",
			}
		}
		err := tx.Model(&model.UserFriend{}).Where("id", requestId).Update("is_accepted", true).Update("accepted_time", tools.NewTimeStamp()).Error
		if err != nil {
			tx.Rollback()
			return false, serializer.Response{
				Status: 500,
				Msg:    "服务器异常！",
			}
		}
		// 检查并处理b向a的反向申请
		var reverseRequest model.UserFriend
		err = tx.Where("friend_id =? AND user_id =? ", userFriend.UserID, userFriend.FriendID).First(&reverseRequest).Error
		if err == nil {
			// 如果找到反向且pending的申请，设为accepted
			err = tx.Model(&model.UserFriend{}).
				Where("id =?", reverseRequest.ID).
				Update("is_accepted", true).Update("accepted_time", tools.NewTimeStamp()).Error
			if err != nil {
				tx.Rollback()
				return false, serializer.Response{
					Status: 500,
					Msg:    "服务器异常！",
				}
			} else {
				tx.Commit()
				return true, serializer.Response{
					Status: 200,
					Msg:    "已接受好友申请!",
				}
			}
		} else {
			if errors.Is(err, gorm.ErrRecordNotFound) {

				tx.Commit()
				return true, serializer.Response{
					Status: 200,
					Msg:    "已接受好友申请!",
				}
			} else {
				tx.Rollback()
				return false, serializer.Response{
					Status: 500,
					Msg:    "服务器异常！",
				}
			}
		}
	} else {
		return false, serializer.Response{
			Status: 400,
			Msg:    "参数异常！",
		}
	}
}

func Friends(id uint) serializer.Response {
	type User struct {
		ID       int    `json:"id"`
		UserName string `json:"userName"`
		Avatar   string `json:"avatar"`
	}
	var friendUserIDs []uint
	// 获取作为fromUserID时的好友ID
	if err := model.DB.Model(&model.UserFriend{}).
		Where("user_id =? AND is_accepted =?", id, true).
		Pluck("friend_id", &friendUserIDs).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return tools.ServiceError()
	}

	var moreFriendUserIDs []uint
	if err := model.DB.Model(&model.UserFriend{}).
		Where("friend_id =? AND is_accepted =?", id, true).
		Pluck("user_id", &moreFriendUserIDs).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return tools.ServiceError()
	}

	friendUserIDs = append(friendUserIDs, moreFriendUserIDs...)

	var friends []User
	// 根据好友ID获取用户信息
	model.DB.Where("id IN?", friendUserIDs).Find(&friends)
	return serializer.Response{
		Status: 200,
		Data:   friends,
		Msg:    "获取成功!",
	}
}
