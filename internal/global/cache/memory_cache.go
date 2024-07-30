package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type MemoryCache struct {
	cacheAdapter *cache.Cache
}

func GetMemoryCache(defaultExpiration, cleanupInterval time.Duration) Cache {
	return Cache(&MemoryCache{cache.New(defaultExpiration, cleanupInterval)})
}

// 加入缓存
func (c *MemoryCache) SetCache(key string, value string, expiration time.Duration) error {
	c.cacheAdapter.Set(key, value, expiration)
	return nil
}

func (c *MemoryCache) GetCache(key string) (string, bool) {
	if value, exist := c.cacheAdapter.Get(key); exist {
		return value.(string), exist
	}
	return "", false
}

// 删除 cache
func (c *MemoryCache) DeleteCache(key string) {
	c.cacheAdapter.Delete(key)
}
