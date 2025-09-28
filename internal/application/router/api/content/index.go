package content

import (
	"ai-software-copyright-server/internal/application/router/api/content/image"
)

type RouterGroup struct {
	Image image.RouterGroup
}

// @title Tool Content-API 文档
// @version v0.0.1
// @description Tool 服务
// @host localhost:8888
// @schemes http https
// @produce json
// @securityDefinitions.apikey content
// @in header
// @name Content-Authorization
// @BasePath /api/content
// @contact.name nineya
// @contact.url https://tool.nineya.com
// @contact.email 361654768@qq.com
var ApiRouterGroup = new(RouterGroup)
