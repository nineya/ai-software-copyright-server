package response

type ZhipuAiChatResponse struct {
	Id      string                  `json:"id" form:"id"`
	Created int64                   `json:"created" form:"created"`
	Model   string                  `json:"model" form:"model"`     //所要调用的模型编码
	Choices []ZhipuAiChatChoiceItem `json:"choices" form:"choices"` //当前对话的模型输出内容
	Usage   ZhipuAiChatUsage        `json:"usage" form:"usage"`     //结束时返回本次模型调用的 tokens 数量统计。
}

type ZhipuAiViewResponse struct {
	Created       int64                          `json:"created" form:"created"`
	Data          []ZhipuAiViewDataItem          `json:"data" form:"data"`
	ContentFilter []ZhipuAiViewContentFilterItem `json:"content_filter" form:"content_filter"` //结束时返回本次模型调用的 tokens 数量统计。
}

type ZhipuAiViewDataItem struct {
	Url string `json:"url" form:"url"` //图片链接。图片的临时链接有效期为 30天，请及时转存图片。
}

type ZhipuAiViewContentFilterItem struct {
	Role  string `json:"role" form:"role"`   //角色
	Level string `json:"level" form:"level"` //严重程度 level 0-3，level 0表示最严重，3表示轻微
}

type ZhipuAiChatChoiceItem struct {
	Index        int                `json:"index" form:"index"`
	FinishReason string             `json:"finish_reason" form:"finish_reason"` //模型推理终止的原因。
	Message      ZhipuAiChatMessage `json:"message" form:"message"`
}

type ZhipuAiChatMessage struct {
	Role    string `json:"role" form:"role"`       //角色
	Content string `json:"content" form:"content"` //聊天内容
}

type ZhipuAiChatUsage struct {
	PromptTokens     int `json:"prompt_tokens" form:"prompt_tokens"`         //用户输入的 tokens 数量
	CompletionTokens int `json:"completion_tokens" form:"completion_tokens"` //模型输出的 tokens 数量
	TotalTokens      int `json:"total_tokens" form:"total_tokens"`           //总 tokens 数量
}
