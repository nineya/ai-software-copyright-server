package task

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/router/socket"
	"ai-software-copyright-server/internal/global"
	"fmt"
)

// socket 心跳
func SocketHeartbeatTask() {
	global.LOG.Info("开始发送Socket心跳")
	socketClients := socket.GetWebSocket().Clients
	for userId, client := range socketClients {
		err := client.WriteMessage(enum.SocketMessageType(1), "收到心跳")
		global.LOG.Info(fmt.Sprintf("发送心跳（%d）V%s:%+v", userId, client.Version, err))
	}
	global.LOG.Info("完成发送Socket心跳")
}
