package user

import (
	"ai-software-copyright-server/internal/application/router/api/user/cdkey"
	"ai-software-copyright-server/internal/application/router/api/user/credits"
	"ai-software-copyright-server/internal/application/router/api/user/netdisk"
	"ai-software-copyright-server/internal/application/router/api/user/public"
	"ai-software-copyright-server/internal/application/router/api/user/qrcode"
	"ai-software-copyright-server/internal/application/router/api/user/redbook"
	"ai-software-copyright-server/internal/application/router/api/user/software_copyright"
	"ai-software-copyright-server/internal/application/router/api/user/study"
	"ai-software-copyright-server/internal/application/router/api/user/user"
)

type RouterGroup struct {
	Cdkey             cdkey.RouterGroup
	Credits           credits.RouterGroup
	Netdisk           netdisk.RouterGroup
	Public            public.RouterGroup
	Qrcode            qrcode.RouterGroup
	Redbook           redbook.RouterGroup
	SoftwareCopyright software_copyright.RouterGroup
	Study             study.RouterGroup
	User              user.RouterGroup
}

// @title Tool User-API 文档
// @version v0.0.1
// @description Tool 服务
// @host localhost:8888
// @schemes http https
// @produce json
// @securityDefinitions.apikey user
// @in header
// @name User-Authorization
// @BasePath /api/user
// @contact.name nineya
// @contact.url https://tool.nineya.com
// @contact.email 361654768@qq.com
var ApiRouterGroup = new(RouterGroup)
