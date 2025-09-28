package response

import "ai-software-copyright-server/internal/application/model/table"

type QrcodeBuildResponse struct {
	UserBuyResponse
	Content string `json:"content"`
}

type QrcodeAddImageResponse struct {
	UserBuyResponse
	table.Qrcode
}
