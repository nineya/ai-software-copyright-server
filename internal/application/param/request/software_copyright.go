package request

import "ai-software-copyright-server/internal/application/model/enum"

type SCTriggerParam struct {
	Id     int64                      `json:"id" form:"id" binding:"required" label:"软著申请ID"`
	Mode   enum.SoftwareCopyrightMode `json:"mode" form:"mode" label:"生成模式"`
	ApiKey *string                    `json:"apiKey" form:"apiKey" label:"ApiKey"`
}
