package service

import (
	"Gin-WebSocket/conf"
	"Gin-WebSocket/pkg/e"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func (manager *ClientManager) Start() {

	for {
		fmt.Println("_________________监听管道通信_____________________")
		select {
		case conn := <-Manager.Register:
			fmt.Printf("有新连接:%v", conn.ID)
			Manager.Clients[conn.ID] = conn //把新连接放进用户管理中
			replyMsg := ReplyMsg{
				Code:    e.WebsocketSuccess,
				Content: "链接到服务器了！",
			}
			msg, _ := json.Marshal(replyMsg)
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
			fmt.Println(Manager.Clients)
		case conn := <-Manager.Unregister:
			fmt.Printf("链接失败:%v", conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok {
				replyMsg := ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "链接中断！",
				}
				msg, _ := json.Marshal(&replyMsg)
				_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
				close(conn.Send)
				delete(Manager.Clients, conn.ID)
			}
		case broadcast := <-Manager.Broadcast: // 1发送给2
			message := broadcast.Message
			sendId := broadcast.Client.SendID //2接收1的消息
			flag := false                     //默认对方不在线
			for id, conn := range Manager.Clients {
				if id != sendId {
					continue
				}
				select {
				case conn.Send <- message:
					flag = true
				default:
					close(conn.Send)
					delete(Manager.Clients, id)
				}
			}
			id := broadcast.Client.ID //1->2
			if flag {
				replyMsg := ReplyMsg{
					Code:    e.WebsocketOnlineReply,
					Content: "对方在线应答！",
				}
				msg, _ := json.Marshal(&replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				err := InsertMsg(conf.MongoDBName, id, string(message), 1, int64(3*month)) //1表示已读
				if err != nil {
					logrus.Info("InsertMsg Err:", err)
				}
			} else {
				replyMsg := ReplyMsg{
					Code:    e.WebsocketOfflineReply,
					Content: "对方不在线！",
				}
				msg, err := json.Marshal(&replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				err = InsertMsg(conf.MongoDBName, id, string(message), 0, int64(3*month))
				if err != nil {
					logrus.Info("InsertMsg Err:", err)
				}
			}
		}
	}
}
