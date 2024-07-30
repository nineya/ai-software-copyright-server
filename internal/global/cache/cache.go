package cache

import "time"

type Cache interface {
	// 设置缓存
	SetCache(key string, value string, expiration time.Duration) error

	// 取得缓存，并通过bool返回是否成功取得缓存
	GetCache(key string) (string, bool)

	// 删除缓存
	DeleteCache(key string)
}
