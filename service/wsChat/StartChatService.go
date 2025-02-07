package wsChat

import (
	"Gin-WebSocket/conf"
	"github.com/sirupsen/logrus"
	"strconv"
)

func StartChatService() {
	for {
		select {
		case conn := <-server.register:
			logrus.Info("用户注册")
			//将新连接放进server中
			server.connections[conn.userInfo.ID] = conn
			logrus.Info("当前在线用户数据：", len(server.connections))
		case conn := <-server.unregister:
			//注销,检查是否存在
			if _, exists := server.connections[conn.userInfo.ID]; exists {
				// 确认连接存在才进行后续操作
				delete(server.connections, conn.userInfo.ID)
				close(conn.send)
				// 将连接设置为nil，防止后续误操作
				conn = nil
				logrus.Info("用户注销")
				logrus.Info("当前在线用户数据：", len(server.connections))
			}
		case msg := <-server.sendToMsg:
			jsonMsg, err := jsonBytesToStruct(msg)
			if err != nil {
				logrus.Info("转换失败！", err)
			}
			//发送给指定用户
			if _, exists := server.connections[jsonMsg.To]; exists {
				//在线
				logrus.Info("在线：", jsonMsg.To)
				server.connections[jsonMsg.To].send <- msg //发送给TO
				if err := InsertMsg(conf.MongoDBName, createID(strconv.Itoa(jsonMsg.From), strconv.Itoa(jsonMsg.To)), string(msg), 1, int64(3*month)); err != nil {
					logrus.Error("MangoDB 插入出错！", err)
				}
			} else {
				//不在线
				logrus.Info("不在线", jsonMsg.To)
				if err := InsertMsg(conf.MongoDBName, createID(strconv.Itoa(jsonMsg.From), strconv.Itoa(jsonMsg.To)), string(msg), 0, int64(3*month)); err != nil {
					logrus.Error("MangoDB 插入出错！", err)
				}
			}
		}
	}
}

//err := InsertMsg(conf.MongoDBName, id, string(message), 1, int64(3*month)) //1表示已读
//const month = 60 * 60 * 24 * 30
//			r1, _ := cache.RedisClient.Get(manager.ID).Result()
//			r2, _ := cache.RedisClient.Get(manager.SendID).Result()
//cache.RedisClient.Incr(manager.ID)
//				//链接建立一个月就会到期
//				_, _ = cache.RedisClient.Expire(manager.ID, time.Hour*24*30).Result() //一个月过期
