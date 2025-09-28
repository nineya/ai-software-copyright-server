package html

import (
	"ai-software-copyright-server/internal/application/router/html/feed"
	"ai-software-copyright-server/internal/application/router/html/netdisk"
	"ai-software-copyright-server/internal/application/router/html/qrcode"
	"ai-software-copyright-server/internal/application/router/html/short_link"
	"github.com/gin-gonic/gin"
)

// 页面路由
func InitHtmlRouter(Router *gin.RouterGroup) {
	// 短链
	Router.GET("s/:alias", short_link.Redirect)
	// 网盘网页
	Router.GET("ptm/:id", netdisk.PcToMobile)
	// 二维码活码
	Router.GET("qrcode/:alias", qrcode.Loose)
}

// 资源文件路由
func InitResourceRouter(Router *gin.RouterGroup) {
	Router.GET("robots.txt", feed.Robots)
	Router.GET("sitemap.xml", feed.SitemapXml)
	Router.GET("sitemap.html", feed.SitemapHtml)
	// 网盘搜索
	Router.GET("netdisk/search", netdisk.Index)
	Router.GET("netdisk/search/index.html", netdisk.Index)
	Router.GET("netdisk/search/search.html", netdisk.Search)
	Router.GET("netdisk/search/detail.html", netdisk.Detail)
	Router.GET("netdisk/search/detail/:id", netdisk.Detail)
}
