package global

import (
	"ai-software-copyright-server/internal/application/model/common"
	"ai-software-copyright-server/internal/config"
	"ai-software-copyright-server/internal/global/cache"
	"ai-software-copyright-server/internal/utils/jwt"
	"embed"
	"github.com/go-pay/gopay/wechat/v3"
	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

var (
	WORK_DIR    string
	CONFIG      config.Config
	JWT         *jwt.JWT
	DB          *xorm.Engine
	REDIS       *redis.Client
	FS          embed.FS
	CACHE       cache.Cache
	LOG         *zap.Logger
	HTML_RENDER common.HtmlRender
	CRON        *cron.Cron
	WECHAT_PAY  *wechat.ClientV3
	SOCKET      common.Socket
)
