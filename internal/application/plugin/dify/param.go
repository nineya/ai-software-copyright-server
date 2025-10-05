package dify

type DifyChatMessageParam struct {
	Query            string                `json:"query" label:"用户输入/提问内容"`
	Inputs           map[string]any        `json:"inputs,omitempty" label:"App 定义的各变量值"`
	ResponseMode     string                `json:"response_mode,omitempty" label:"响应模式"` // streaming 流式模式，blocking 阻塞模式
	User             string                `json:"user,omitempty" label:"用户标识"`
	ConversationId   string                `json:"conversation_id,omitempty" label:"会话 ID"` // 放空则创建新会话
	Files            []DifyChatMessageFile `json:"files,omitempty" label:"文件列表"`
	AutoGenerateName bool                  `json:"auto_generate_name,omitempty" label:"自动生成标题"` // 默认 true
}

type DifyChatMessageFile struct {
	Type           string `json:"type,omitempty" form:"type" label:"文件类型"`
	TransferMethod string `json:"transfer_method" form:"transfer_method" label:"传递方式"`            // remote_url 文件地址，local_file 上传文件
	Url            string `json:"url,omitempty" form:"url" label:"文件地址"`                          // 仅当传递方式为 remote_url 时
	UploadFileId   string `json:"upload_file_id,omitempty" form:"upload_file_id" label:"上传文件 ID"` // 仅当传递方式为 local_file 时
}
