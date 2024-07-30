package initialize

import (
	"github.com/pkg/errors"
	"time"
	"tool-server/internal/global"
	"tool-server/internal/global/cache"
)

func InitCache() {
	cacheType := global.CONFIG.Server.Cache
	switch cacheType {
	case "redis":
		if global.REDIS == nil {
			panic(errors.New("Redis is not initialized"))
		}
		global.CACHE = cache.GetRedisCache(global.REDIS)
		break
	default:
		global.CACHE = cache.GetMemoryCache(5*time.Minute, time.Minute)
	}
}
