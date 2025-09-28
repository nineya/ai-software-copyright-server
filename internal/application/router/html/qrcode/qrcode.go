package qrcode

import (
	"ai-software-copyright-server/internal/application/param/response"
	qrcodeSev "ai-software-copyright-server/internal/application/service/qrcode"
	"ai-software-copyright-server/internal/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
)

// 活码解析
func Loose(c *gin.Context) {
	htmlResponse := response.GenerateHtmlResult(c)
	alias := c.Param("alias")
	targetUrl := ""
	if alias != "" {
		qrcode, err := qrcodeSev.GetQrcodeService().GetByAlias(alias)
		if err != nil {
			global.LOG.Error(fmt.Sprintf("查询活链失败: %+v", err))
		} else if qrcode.Id > 0 {
			_ = qrcodeSev.GetQrcodeService().UpdateVisitsIncreaseById(qrcode.Id, c.Request.UserAgent())
			if qrcode.TargetUrls != nil && len(qrcode.TargetUrls) > 0 {
				targetUrl = qrcode.TargetUrls[rand.Intn(len(qrcode.TargetUrls))]
			}
		}
	}
	htmlResponse.OkWithData("feature/qrcode/loose.html", gin.H{
		"Alias":     alias,
		"TargetUrl": targetUrl,
	})
}
