package html

import (
	"ai-software-copyright-server/internal/application/router/html/feed"
	"github.com/gin-gonic/gin"
)

// 页面路由
func InitHtmlRouter(Router *gin.RouterGroup) {
}

// 资源文件路由
func InitResourceRouter(Router *gin.RouterGroup) {
	Router.GET("robots.txt", feed.Robots)
	Router.GET("download/softwareCopyright/:id", feed.SoftwareCopyright)
}
