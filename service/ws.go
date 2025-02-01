package service

import (
	"Gin-WebSocket/cache"
	"Gin-WebSocket/conf"
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

// 发送结构体
type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

// 回复消息结构体
type ReplyMsg struct {
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

// 用户结构体
type Client struct {
	ID     string
	SendID string
	Socket *websocket.Conn
	Send   chan []byte
}

// 广播类（包括广播内容与源用户）
type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

// 用户管理
type ClientManager struct {
	Clients    map[string]*Client //存储所有当前连接到服务器的客户端
	Broadcast  chan *Broadcast    //用于广播消息给所有已连接的客户端
	Reply      chan *Client
	Register   chan *Client //注册
	Unregister chan *Client //注销
}

// 序列化
/*
omitempty:
去除不必要的零值字段
确保只有有意义的数据被返回
*/
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var Manager = ClientManager{
	Clients:    make(map[string]*Client),
	Broadcast:  make(chan *Broadcast),
	Register:   make(chan *Client),
	Reply:      make(chan *Client),
	Unregister: make(chan *Client),
}

func CreateID(uid, toUid string) string {
	return uid + "->" + toUid
}

func Handler(ctx *gin.Context) {
	// 获取请求参数中的uid
	uid := ctx.Query("uid")
	// 获取请求参数中的toUid
	toUid := ctx.Query("toUid")

	// 创建websocket升级器，并设置跨域策略为允许所有源
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(ctx.Writer, ctx.Request, nil) // 升级websocket协议

	// 如果升级失败，返回404错误
	if err != nil {
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}

	// 创建用户实例
	client := &Client{
		ID:     CreateID(uid, toUid), // uid发送给toUid的标识
		SendID: CreateID(toUid, uid), // toUid接收到uid的标识
		Socket: conn,
		Send:   make(chan []byte), // 创建发送消息的通道
	}

	// 用户注册到用户管理上
	// 将用户实例发送到用户管理器的Register通道中
	Manager.Register <- client

	// 在后台goroutine中读取消息
	go client.Read()

	// 在后台goroutine中写入消息
	go client.Write()
}

func (manager *Client) Read() {
	defer func() {
		//将用户进行注销操作
		Manager.Unregister <- manager
		_ = manager.Socket.Close() //关闭该socket
	}()
	for {
		manager.Socket.PongHandler()
		sendMsg := new(SendMsg)
		//client.Socket.ReadMessage()//字符串类型
		if err := manager.Socket.ReadJSON(&sendMsg); err != nil { //json类型
			logrus.Info("数据格式不正确！", err)
			Manager.Unregister <- manager //注销链接
			_ = manager.Socket.Close()    //关闭socket
			break                         //关闭循环
		}
		if sendMsg.Type == 1 { //发送消息
			//在缓存中找 1->2 和 2->1
			r1, _ := cache.RedisClient.Get(manager.ID).Result()
			r2, _ := cache.RedisClient.Get(manager.SendID).Result()
			if r1 > "3" && r2 == "" { //(1->2)发三条以上，2不回，则销毁资源
				replyMsg := ReplyMsg{
					Code:    e.WebsocketOfflineReply,
					Content: e.GetMsg(e.WebsocketLimit),
				}
				msg, _ := json.Marshal(&replyMsg)                           //序列化
				_ = manager.Socket.WriteMessage(websocket.TextMessage, msg) //使用空标识符 _ 来忽略 WriteMessage 方法返回的错误。这意味着调用 client.Socket.WriteMessage 方法时，即使出现错误，程序也不会对该错误做任何处理，会继续往下执行。
				//err := client.Socket.WriteMessage(websocket.TextMessage, msg) //通过这种方式，我们可以在后续代码中对错误进行检查和处理。正确处理错误能让程序更加健壮和可靠。
				continue
			} else {
				/*
					Incr 函数会先将键的值初始化为 1，然后返回这个初始值。
					如果键已经存在且它的值是一个数字，则会将该值增加 1，并返回增加后的值。
				*/
				cache.RedisClient.Incr(manager.ID)
				//链接建立一个月就会到期
				_, _ = cache.RedisClient.Expire(manager.ID, time.Hour*24*30).Result() //一个月过期
			}
			//广播该消息
			Manager.Broadcast <- &Broadcast{
				Client:  manager,
				Message: []byte(sendMsg.Content), //消息
			}
		} else if sendMsg.Type == 2 {
			//获取历史消息
			timeT, err := strconv.Atoi(sendMsg.Content) //string to int
			if err != nil {
				timeT = 999999999
			}
			fmt.Println("SendID:", manager.SendID, "ID:", manager.ID)
			results, _ := FindMany(conf.MongoDBName, manager.SendID, manager.ID, int64(timeT), 10) //获取10条历史消息

			if len(results) > 10 {
				results = results[:10]
			} else if len(results) == 0 {
				replyMsg := ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "没有更多历史记录了！",
				}
				msg, _ := json.Marshal(&replyMsg)
				_ = manager.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			}
			for _, result := range results {
				replyMsg := ReplyMsg{
					From:    result.From,
					Content: result.Msg,
				}
				msg, _ := json.Marshal(&replyMsg)
				_ = manager.Socket.WriteMessage(websocket.TextMessage, msg)
			}

		}

	}

}

func (client *Client) Write() {
	defer func() {
		_ = client.Socket.Close()
	}()
	for {
		select {
		//取消息
		case message, ok := <-client.Send:
			if !ok {
				_ = client.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			replyMsg := ReplyMsg{
				Code:    e.WebsocketSuccess,
				Content: fmt.Sprintf("%s", string(message)),
			}
			msg, _ := json.Marshal(&replyMsg)
			_ = client.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
