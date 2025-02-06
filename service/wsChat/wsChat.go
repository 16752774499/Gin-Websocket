package wsChat

import (
	"Gin-WebSocket/conf"
	"Gin-WebSocket/model"
	"Gin-WebSocket/pkg/e"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"sync"
)

// 消息过期时长
const month = 60 * 60 * 24 * 30

//用户在线状态管理

type Message struct {
	Type    int    `json:"type"`
	From    int    `json:"from"` //发送者
	To      int    `json:"to"`   //接收者
	Content string `json:"content"`
	Other   []byte `json:"other,omitempty"`
}

type UserInfo struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Avatar   string `json:"avatar"`
}

type ChatUserConnManager struct {
	User     map[int]*websocket.Conn
	UserInfo map[int]*UserInfo
	mu       sync.RWMutex
}

var chatUserConnManager = &ChatUserConnManager{
	User:     make(map[int]*websocket.Conn),
	UserInfo: make(map[int]*UserInfo),
}

func Chat(ctx *gin.Context) {

	userId, _ := strconv.Atoi(ctx.Query("uid"))
	//升级连接
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { //允许跨域
			return true
		},
	}).Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		//直接通过http包返回404，并终止这次操作
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}
	chatUserConnManager.AddConn(userId, conn)
	defer chatUserConnManager.DelConn(userId)
	//直接监听消息(只关注消息转发)，用户管理有隔壁heartbeat做

	for {
		var message Message
		t, p, err := conn.ReadMessage()
		if err != nil {
			logrus.Error("Heartbeat read error:", zap.Any("错误", err))
			break
		}
		fmt.Println("消息类型:", t, "消息内容:", string(p), "错误:", err)
		if err := json.Unmarshal(p, &message); err != nil {
			fmt.Printf("JSON 解析错误: %v，原始消息: %s\n", err, string(p))
		}
		//判断

		switch message.Type {
		case e.ChatSystemMsg:
			fmt.Println("系统消息！")
		case e.ChatUserCommonMsg:
			fmt.Println("普通消息！")
			SendMsg(userId, message)
		default:
			fmt.Println("当前在线用户数")
			fmt.Println(chatUserConnManager.UserInfo)

		}
		fmt.Println("当前在线用户:", chatUserConnManager.UserInfo)
	}

}

func (chat *ChatUserConnManager) AddConn(userId int, conn *websocket.Conn) {
	chat.mu.Lock()
	defer chat.mu.Unlock()
	chat.User[userId] = conn
	//添加用户信息
	chat.UserInfo[userId] = getUserInfo(userId)
	conf.Log.Info("用户上线聊天", zap.Int("用户id", userId))
}
func (chat *ChatUserConnManager) DelConn(userId int) {
	chat.mu.Lock()
	defer chat.mu.Unlock()
	delete(chat.User, userId)
	delete(chat.UserInfo, userId)
	conf.Log.Info("聊天用户下线", zap.Int("用户id", userId))
}

func getUserInfo(userId int) *UserInfo {
	userInfo := &UserInfo{}
	model.DB.Model(&model.User{}).Where("id = ?", userId).First(&userInfo)
	return userInfo
}

func getChatId(fromId, toId int) string {
	return fmt.Sprintf("%d->%d", fromId, toId)
}

func (c *ChatUserConnManager) isUserOnline(id int) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	// 检查用户连接是否存在
	if c.User[id] != nil {
		return true
	}
	return false
}
