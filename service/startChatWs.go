package service

import (
	"Gin-WebSocket/conf"
	"Gin-WebSocket/pkg/e"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

func (manager *ClientManager) Start() {

	for {
		fmt.Println("_________________监听管道通信_____________________")
		select {
		case conn := <-Manager.Register:
			// 处理新连接注册
			fmt.Printf("有新连接:%v", conn.ID)
			Manager.Clients[conn.ID] = conn // 把新连接放进用户管理中
			// 构造并发送回复消息
			replyMsg := ReplyMsg{
				Code:    e.WebsocketLinkSuccess,
				Content: "链接到服务器了！",
			}
			msg, _ := json.Marshal(replyMsg)
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
			fmt.Println(Manager.Clients)
		case conn := <-Manager.Unregister:
			// 处理连接注销
			fmt.Printf("链接失败:%v", conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok {
				// 构造并发送回复消息
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
			// 处理广播消息
			message := broadcast.Message
			sendId := broadcast.Client.SendID // 2接收1的消息
			id := broadcast.Client.ID         // 1->2
			toId, _ := strconv.Atoi(strings.Split(sendId, "-")[0])

			//检查对方客户端是否在线
			cflag := onlineManager.isUserOnline(toId)

			//socket连接是否在线
			sflag := false // 默认对方不在线 （反向的socket连接是否在线）
			// 遍历所有连接，寻找目标连接
			for id, conn := range Manager.Clients {
				fmt.Println(id, conn.ID, sendId)
				if id != sendId {
					continue
				}
				select {
				case conn.Send <- message:
					sflag = true
				default:
					close(conn.Send)
					delete(Manager.Clients, id)
				}
			}
			fmt.Printf("用户id为%d的客户端在线状态%t\n", toId, cflag)
			conf.Log.Info("用户id为toId的客户端在线状态cflag\n", zap.Int("toId", toId), zap.Bool("cflag", cflag))
			//回复消息的socket
			fmt.Printf("%s的websocket连接在线状态%t\n", sendId, sflag)
			conf.Log.Info("toId的websocket连接在线状态cflag\n", zap.Int("toId", toId), zap.Bool("cflag", cflag))

			if sflag {
				//socket连接在线，直接发送
				// 对方在线，构造并发送回复消息
				replyMsg := ReplyMsg{
					Code:    e.WebsocketOnlineReply,
					Content: "对方在线应答！",
				}
				msg, _ := json.Marshal(&replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				// 将消息插入数据库，标记为已读
				err := InsertMsg(conf.MongoDBName, id, string(message), 1, int64(3*month)) // 1表示已读
				if err != nil {

					conf.Log.Warn("InsertMsg Err:", zap.Any("err", err))
				}
			} else {
				//客户端离线 - socket连接一定是没有的
				//消息保存，待客户端上线推送
				// 对方不在线，构造并发送回复消息
				replyMsg := ReplyMsg{
					Code:    e.WebsocketOfflineReply,
					Content: "对方不在线！",
				}
				msg, err := json.Marshal(&replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				// 将消息插入数据库，标记为未读
				err = InsertMsg(conf.MongoDBName, id, string(message), 0, int64(3*month))
				if err != nil {
					conf.Log.Warn("InsertMsg Err:", zap.Any("err", err))
				}
			}
		}
	}
}
