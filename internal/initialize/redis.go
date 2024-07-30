package initialize

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"tool-server/internal/global"
)

func InitRedis() {
	if global.CONFIG.Server.Cache != "redis" {
		return
	}
	redisCfg := global.CONFIG.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Host,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DB,       // use default DB
	})
	result, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(errors.Wrap(err, "Redis connection failed"))
	} else {
		global.LOG.Info("Redis connect ping response: " + result)
		global.REDIS = client
	}
}
