package service

import (
	"Gin-WebSocket/conf"
	"fmt"
	"go.uber.org/zap"
	"time"
)

// 检查心跳超时
func checkHeartbeats() {
	// 加锁以保护对onlineManager的并发访问
	onlineManager.mu.Lock()
	defer onlineManager.mu.Unlock() // 解锁操作在函数结束时执行

	now := time.Now()          // 获取当前时间
	timeout := 2 * time.Minute // 定义超时时间：两分钟没心跳认为离线

	// 遍历LastPing映射，检查每个用户的最后ping时间
	for userId, lastPing := range onlineManager.LastPing {
		// 如果当前时间与最后ping时间的差值超过超时时间
		if now.Sub(lastPing) > timeout {
			// 移除超时conn
			// 尝试从Connections映射中获取对应userId的连接
			if conn := onlineManager.Connections[userId]; conn != nil {
				// 关闭连接
				err := conn.Close()
				if err != nil {
					// 如果关闭连接时出错，则直接返回
					return
				}
				// 从Connections映射中删除该userId的连接
				delete(onlineManager.Connections, userId)
				// 从LastPing映射中删除该userId的最后ping时间
				delete(onlineManager.LastPing, userId)
				// 从OnlineUsers映射中删除该userId的在线状态
				delete(onlineManager.OnlineUsers, userId)
				// 广播用户状态变化
				// 广播userId的离线状态
				broadcastUserStatus(userId, false)
			}
		}
	}
	fmt.Printf("当前存活用户数:%d\n", len(onlineManager.Connections))
	conf.Log.Info("当前存活用户数:", zap.Int("OnlineUsersLen", len(onlineManager.LastPing)))
	fmt.Println("用户:", onlineManager.OnlineUsers)
	conf.Log.Info("用户:", zap.Any("OnlineUsers", onlineManager.OnlineUsers))
}

// 心条检测启动程序
func StartHeartbeats() {
	go func() {
		for {
			time.Sleep(time.Minute) //每分钟检查一次
			checkHeartbeats()
		}
	}()
}
