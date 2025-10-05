package dify

type DifyChatMessageResponse struct {
	Event          string `json:"event"`           // 事件类型，固定为 message
	TaskId         string `json:"task_id"`         // 任务 ID，用于请求跟踪和下方的停止响应接口
	Id             string `json:"id"`              // 唯一ID
	MessageId      string `json:"message_id"`      // 消息唯一 ID
	ConversationId string `json:"conversation_id"` // 会话 ID
	Mode           string `json:"mode"`            // App 模式，固定为 chat
	Answer         string `json:"answer"`          // 完整回复内容
	CreatedAt      int    `json:"created_at"`      // 消息创建时间戳
}

type DifyChatMessageSSEResponse struct {
	Event          string `json:"event"`           // 事件类型，固定为 message
	TaskId         string `json:"task_id"`         // 任务 ID，用于请求跟踪和下方的停止响应接口
	Id             string `json:"id"`              // 唯一ID
	MessageId      string `json:"message_id"`      // 消息唯一 ID
	ConversationId string `json:"conversation_id"` // 会话 ID
	Mode           string `json:"mode"`            // App 模式，固定为 chat
	Answer         string `json:"answer"`          // 完整回复内容
	CreatedAt      int    `json:"created_at"`      // 消息创建时间戳
}
