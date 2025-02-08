package wsChat

import "C"
import (
	"Gin-WebSocket/conf"
	"Gin-WebSocket/model"
	"Gin-WebSocket/pkg/e"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

// 消息过期时长
const month = 60 * 60 * 24 * 30

// 历史消息
type MessageHistory struct {
}

//用户在线状态管理

type Message struct {
	Type    int    `json:"type"`
	From    int    `json:"from"` //发送者
	To      int    `json:"to"`   //接收者
	Content string `json:"content"`
	Time    int64  `json:"time"`
	Other   []byte `json:"other,omitempty"`
}
type UserInfo struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Avatar   string `json:"avatar"`
}
type Connection struct {
	conn     *websocket.Conn
	send     chan []byte
	userInfo *UserInfo
}
type Server struct {
	connections map[int]*Connection
	register    chan *Connection //注册
	unregister  chan *Connection //注销
	sendToMsg   chan []byte      //发送消息信号
}

var server = &Server{
	connections: make(map[int]*Connection),
	register:    make(chan *Connection),
	unregister:  make(chan *Connection),
	sendToMsg:   make(chan []byte),
}

// 处理客户端链接
func HandleChat(ctx *gin.Context) {

	userId, _ := strconv.Atoi(ctx.Query("uid"))
	logrus.Info("uid:", userId)
	//升级为websocket链接
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true //允许跨域
		},
	}).Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logrus.Info("websocket upgrade fail", err)
	}
	//创建链接对象
	connection := &Connection{
		conn:     conn,
		send:     make(chan []byte, 256),
		userInfo: getUserInfo(userId),
	}
	//注册进server
	server.register <- connection
	//给每个用户创建读写携程
	go connection.Read()
	go connection.Write()

}

func (client *Connection) Read() {
	//注销
	defer func() {
		server.unregister <- client
	}()
	/*
			在 WebSocket 协议中，Ping 是一种控制帧 ，用于检查对等方（服务器或客户端）是否仍然连接并且正常运行。PingHandler 是 websocket.Conn 类型的一个方法，用于获取或设置用于处理传入 Ping 帧的处理函数。
		如果调用不带参数的 client.conn.PingHandler()，它的作用是返回当前用于处理 Ping 帧的处理函数。如果调用 client.conn.PingHandler(func(messageType int, p []byte) error) 并传入一个函数，那么这个传入的函数就会被设置为新的 Ping 帧处理函数。
		默认情况下，当 websocket.Conn 接收到一个 Ping 帧时，会回复一个 Pong 帧
	*/
	client.conn.PingHandler()

	for {
		client.conn.PongHandler()
		message := new(Message)
		if err := client.conn.ReadJSON(&message); err != nil {
			logrus.Info("read fail：", err)
			server.unregister <- client
			break
		}

		//读取消息
		switch message.Type {
		case e.ChatSystemMsg:
			logrus.Info("系统消息", message)
		case e.ChatUserCommonMsg:
			bytes, err := structToJsonBytes(message)
			if err != nil {
				return
			}
			server.sendToMsg <- bytes
		case e.ChatMessageFile:
			bytes, err := structToJsonBytes(message)
			if err != nil {
				return
			}
			server.sendToMsg <- bytes
		case e.ChatMessageHistory:
			//历史消息
			client.ChatMessageHistory(*message)
		case e.ChatMessageNew:
			client.ChatMessageNew(*message)

		default:
			logrus.Info("其他")
		}

	}

}

func (client *Connection) Write() {
	// 延迟执行函数，确保在函数退出时执行
	defer func() {
		server.unregister <- client
	}()

	// 循环处理消息发送
	for {
		select {
		// 从发送通道接收消息
		case msg, ok := <-client.send:
			if !ok {
				// 如果通道已关闭，发送关闭消息并返回
				_ = client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// 将消息转换为JSON格式
			writeMsg, _ := jsonBytesToStruct(msg)
			// 发送文本消息
			_ = client.conn.WriteJSON(writeMsg)
			logrus.Info("转发了！")

		}
	}
}

func getUserInfo(userID int) *UserInfo {
	userINfo := UserInfo{}
	model.DB.Model(model.User{}).Where("id =?", userID).Find(&userINfo)
	return &userINfo
}

// 结构体转json二进制
func structToJsonBytes(bytes *Message) ([]byte, error) {
	return json.Marshal(bytes)
}

// json二进制转结构体
func jsonBytesToStruct(jsonBytes []byte) (Message, error) {
	var msg Message
	err := json.Unmarshal(jsonBytes, &msg)
	return msg, err
}

func createID(fromId, toId string) string {
	return fromId + "->" + toId
}

func (client *Connection) ChatMessageHistory(message Message) {
	///获取历史消息

	timeT, err := strconv.Atoi(message.Content)
	if err != nil {
		timeT = int(time.Now().Unix())
	}

	fmt.Println("SendID:", strconv.Itoa(message.To), strconv.Itoa(message.From))
	results, err := FindMany(conf.MongoDBName, createID(strconv.Itoa(message.To), strconv.Itoa(message.From)), createID(strconv.Itoa(message.From), strconv.Itoa(message.To)), int64(timeT), 10) //获取10条历史消息
	if err != nil {
		replyMsg := Message{
			Type:    e.ChatMessageHistory,
			Content: "err",
			From:    0,
			To:      message.From,
			Time:    time.Now().Unix(),
		}
		client.conn.WriteJSON(replyMsg)
		return
	}
	logrus.Info(len(results), results)
	if len(results) > 10 {
		results = results[:10]
	} else if len(results) == 0 {
		replyMsg := Message{
			Type:    e.ChatMessageHistory,
			Content: "null",
			From:    0,
			To:      message.From,
			Time:    time.Now().Unix(),
		}

		client.conn.WriteJSON(replyMsg)
		return
	}
	data, _ := json.Marshal(results)
	replyMsg := Message{
		Type:    e.ChatMessageHistory,
		Content: string(data),
		From:    0,
		To:      message.From,
		Time:    time.Now().Unix(),
	}
	//messageHistory, _ := structToJsonBytes(&replyMsg)

	client.conn.WriteJSON(replyMsg)

}

// 有空再写
func (client *Connection) ChatMessageNew(message Message) {

}

//1738939184901
//1738942027
