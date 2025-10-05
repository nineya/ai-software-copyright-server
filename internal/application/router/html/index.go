package html

import (
	"ai-software-copyright-server/internal/application/router/html/feed"
	"ai-software-copyright-server/internal/application/router/html/qrcode"
	"ai-software-copyright-server/internal/application/router/html/short_link"
	"github.com/gin-gonic/gin"
)

// 页面路由
func InitHtmlRouter(Router *gin.RouterGroup) {
	// 短链
	Router.GET("s/:alias", short_link.Redirect)
	// 二维码活码
	Router.GET("qrcode/:alias", qrcode.Loose)
}

// 资源文件路由
func InitResourceRouter(Router *gin.RouterGroup) {
	Router.GET("robots.txt", feed.Robots)
}
