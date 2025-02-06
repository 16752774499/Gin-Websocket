package wsHeartbeat

import (
	"Gin-WebSocket/conf"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// 客户端在线状态管理
type OnlineManager struct {
	// OnlineUsers 存储在线用户的映射关系，其中键是用户ID（int类型），值是用户是否在线的标记（bool类型）。
	// 用户ID (userId) 是唯一的标识符，而在线标记 (isOnline) 表示该用户当前是否在线。
	OnlineUsers map[int]bool // userId -> isOnline

	// Connections 存储每个用户的心跳连接，其中键是用户ID（int类型），值是指向该用户websocket连接的指针（*websocket.Conn）。
	// 这允许系统通过用户ID直接访问并管理用户的心跳连接。
	Connections map[int]*websocket.Conn // userId -> heartbeat connection

	// LastPing 存储每个用户最后一次发送ping消息的时间，其中键是用户ID（int类型），值是最后一次ping的时间戳（time.Time类型）。
	// 这用于跟踪用户的活跃状态，以及检测用户是否已断开连接。
	LastPing map[int]time.Time // userId -> last ping time

	// mu 是一个读写互斥锁（sync.RWMutex），用于确保对OnlineUsers、Connections和LastPing的并发访问是安全的。
	// 这在多线程或并发环境下尤其重要，以避免数据竞争和不一致性问题。
	mu sync.RWMutex
}

// 初始化在线管理器
var onlineManager = &OnlineManager{
	OnlineUsers: make(map[int]bool),
	Connections: make(map[int]*websocket.Conn),
	LastPing:    make(map[int]time.Time),
}

// 在ws路由处注册
func HandleHeartbeat(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Query("uid"))
	//升级链接
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}, // 升级websocket协议
	}).Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		//直接通过http包返回404，并终止这次操作
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}
	//将该用户注册到管理器中
	onlineManager.AddConnection(userId, conn)
	//退出时移除此连接
	defer onlineManager.RemoveConnection(userId)
	//循环监听心跳
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			// 如果是正常关闭连接的错误，直接退出循环
			//if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			//	log.Printf("unexpected close error: %v", err)
			//}
			//break
			conf.Log.Error("Heartbeat read error:", zap.Any("错误", err))
			break
		}
		var heartbeat struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(msg, &heartbeat); err != nil {
			//解析错误，等待下一次
			continue
		}
		//客户端发送ping 服务器相应pong
		if heartbeat.Type == "ping" {
			//更新心跳时间
			onlineManager.UpdateLastPing(userId)

			//相应客户端
			err := conn.WriteJSON(map[string]string{"type": "pong"})
			if err != nil {
				return
			}
		}

	}

}

func (om *OnlineManager) AddConnection(userId int, conn *websocket.Conn) {
	//上锁
	om.mu.Lock()
	//执行完时解锁
	defer om.mu.Unlock()
	om.OnlineUsers[userId] = true
	om.Connections[userId] = conn
	om.LastPing[userId] = time.Now()

	//广播用户上线状态
	broadcastUserStatus(userId, true)
	conf.Log.Info("用户上线:", zap.Int("上线id", userId))
}

func (om *OnlineManager) RemoveConnection(userId int) {
	//上锁
	om.mu.Lock()
	//执行完时解锁
	defer om.mu.Unlock()
	delete(om.Connections, userId)
	delete(om.LastPing, userId)
	delete(om.OnlineUsers, userId)
	//广播下线
	broadcastUserStatus(userId, false)
	conf.Log.Info("用户下线:", zap.Int("下线id", userId))

}

func (om *OnlineManager) UpdateLastPing(userId int) {
	om.mu.Lock()
	defer om.mu.Unlock()
	om.LastPing[userId] = time.Now()
}

// 对外使用，检查用户是否在线
func (om *OnlineManager) isUserOnline(userId int) bool {
	om.mu.RLock()
	defer om.mu.RUnlock()
	return om.OnlineUsers[userId]
}

func broadcastUserStatus(userId int, online bool) {
	status := map[string]interface{}{
		"type":      "user_status",
		"userId":    userId,
		"online":    online,
		"timestamp": time.Now(),
	}
	//广播给所有用户
	for _, conn := range onlineManager.Connections {
		if conn != nil {
			err := conn.WriteJSON(status)
			if err != nil {
				return
			}
		}
	}

}
