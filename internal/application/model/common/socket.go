package common

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type Socket interface {
	// 创建Socket
	NewHandler(c *gin.Context)
	// 发送一个消息
	SendMessage(userId int64, param SocketMessage, result any) (bool, error)
	// 取得客户端
	GetClient(userId int64) *SocketClient
}

// socket消息实体
type SocketMessage struct {
	TaskId     string                 `json:"taskId" form:"taskId" label:"任务ID"`
	NeedResult bool                   `json:"needResult" form:"needResult" label:"是否等待结果"` // 是否需要等待结果
	Timeout    time.Duration          `json:"timeout" form:"timeout" label:"超时时间"`
	Type       enum.SocketMessageType `json:"type" form:"type" binding:"required" label:"消息类型"`
	Data       string                 `json:"data" form:"data" binding:"required" label:"数据内容"`
}

// socket客户端实体
type SocketClient struct {
	sync.Mutex
	Conn    *websocket.Conn `json:"conn"`
	Version string          `json:"version"`
}

func (s *SocketClient) WriteMessage(typ enum.SocketMessageType, data any) error {
	serializedData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return s.WriteJSON(SocketMessage{
		Type: typ,
		Data: string(serializedData),
	})
}

// 发送一个响应，避免并发请求错误，必须绑定指针避免并发错误
func (s *SocketClient) WriteJSON(v interface{}) error {
	s.Lock()
	defer s.Unlock()
	return s.Conn.WriteJSON(v)
}
