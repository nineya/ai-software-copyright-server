package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"tool-server/internal/config"
	"tool-server/internal/global/cache"
	"tool-server/internal/utils/jwt"
	"xorm.io/xorm"
)

var (
	WORK_DIR string
	CONFIG   config.Config
	JWT      *jwt.JWT
	DB       *xorm.Engine
	REDIS    *redis.Client
	CACHE    cache.Cache
	LOG      *zap.Logger
	CRON     *cron.Cron
)
