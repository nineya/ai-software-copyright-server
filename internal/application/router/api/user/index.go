package user

import (
	"ai-software-copyright-server/internal/application/router/api/user/ai"
	"ai-software-copyright-server/internal/application/router/api/user/cdkey"
	"ai-software-copyright-server/internal/application/router/api/user/credits"
	"ai-software-copyright-server/internal/application/router/api/user/flash_picture"
	"ai-software-copyright-server/internal/application/router/api/user/mp"
	"ai-software-copyright-server/internal/application/router/api/user/netdisk"
	"ai-software-copyright-server/internal/application/router/api/user/public"
	"ai-software-copyright-server/internal/application/router/api/user/qrcode"
	"ai-software-copyright-server/internal/application/router/api/user/redbook"
	"ai-software-copyright-server/internal/application/router/api/user/short_link"
	"ai-software-copyright-server/internal/application/router/api/user/study"
	"ai-software-copyright-server/internal/application/router/api/user/time_clock"
	"ai-software-copyright-server/internal/application/router/api/user/user"
)

type RouterGroup struct {
	Ai           ai.RouterGroup
	Cdkey        cdkey.RouterGroup
	Credits      credits.RouterGroup
	FlashPicture flash_picture.RouterGroup
	Mp           mp.RouterGroup
	Netdisk      netdisk.RouterGroup
	Public       public.RouterGroup
	Qrcode       qrcode.RouterGroup
	Redbook      redbook.RouterGroup
	ShortLink    short_link.RouterGroup
	Study        study.RouterGroup
	TimeClock    time_clock.RouterGroup
	User         user.RouterGroup
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
