package response

import (
	"ai-software-copyright-server/internal/application/model/table"
	"github.com/go-pay/gopay/wechat/v3"
)

type CreditsCreateOrderResponse struct {
	*table.CreditsOrder
	AppletParams *wechat.AppletParams `json:"appletParams"`
}
