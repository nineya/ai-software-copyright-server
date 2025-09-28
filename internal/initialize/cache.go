package initialize

import (
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/global/cache"
	"github.com/pkg/errors"
	"time"
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
