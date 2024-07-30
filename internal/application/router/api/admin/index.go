package admin

import (
	"tool-server/internal/application/router/api/admin/admin"
	"tool-server/internal/application/router/api/admin/public"
	"tool-server/internal/application/router/api/admin/redbook"
)

type RouterGroup struct {
	Admin   admin.RouterGroup
	Public  public.RouterGroup
	Redbook redbook.RouterGroup
}

// @title Aurora Admin-API 文档
// @version v0.0.1
// @description Aurora 建站
// @contact.name nineya
// @contact.url https://www.nineya.com
// @contact.email 361654768@qq.com
// @schemes http https
// @host localhost:8888
// @BasePath /api/admin
// @produce json
// @securityDefinitions.apikey admin
// @in header
// @name Admin-Authorization
var ApiRouterGroup = new(RouterGroup)
