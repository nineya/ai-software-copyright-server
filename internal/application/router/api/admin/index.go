package admin

import (
	"ai-software-copyright-server/internal/application/router/api/admin/admin"
	"ai-software-copyright-server/internal/application/router/api/admin/cdkey"
	"ai-software-copyright-server/internal/application/router/api/admin/public"
	"ai-software-copyright-server/internal/application/router/api/admin/redbook"
	"ai-software-copyright-server/internal/application/router/api/admin/user"
)

type RouterGroup struct {
	Admin   admin.RouterGroup
	Cdkey   cdkey.RouterGroup
	Public  public.RouterGroup
	Redbook redbook.RouterGroup
	User    user.RouterGroup
}

// @title Tool Admin-API 文档
// @version v0.0.1
// @description Tool 服务
// @contact.name nineya
// @contact.url https://tool.nineya.com
// @contact.email 361654768@qq.com
// @schemes http https
// @host localhost:8888
// @BasePath /api/admin
// @produce json
// @securityDefinitions.apikey admin
// @in header
// @name Admin-Authorization
var ApiRouterGroup = new(RouterGroup)
