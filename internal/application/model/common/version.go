package common

type Version struct {
	Version     string `json:"version" form:"description"`                         // 版本
	DownloadUrl string `json:"downloadUrl,omitempty" form:"downloadUrl,omitempty"` // 下载地址
	Description string `json:"description,omitempty" form:"description,omitempty"` // 用于描述函数功能。模型会根据这段描述决定函数调用方式。
}
