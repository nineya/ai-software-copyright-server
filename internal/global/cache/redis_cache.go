package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisCache struct {
	redisClient *redis.Client
}

func GetRedisCache(redisClient *redis.Client) Cache {
	return Cache(&RedisCache{redisClient})
}

func (c *RedisCache) SetCache(key string, value string, expiration time.Duration) error {
	return c.redisClient.Set(context.Background(), key, value, expiration).Err()
}

func (c *RedisCache) GetCache(key string) (string, bool) {
	if value, err := c.redisClient.Get(context.Background(), key).Result(); err == nil {
		return value, true
	} else {
		return value, false
	}
}

// 删除 cache
func (c *RedisCache) DeleteCache(key string) {
	c.redisClient.Del(context.Background(), key)
}
