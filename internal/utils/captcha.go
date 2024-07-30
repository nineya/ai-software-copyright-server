package utils

import (
	"go.uber.org/zap"
	"time"
	"tool-server/internal/global"
)

func NewCaptchaStore() *CaptchaStore {
	return &CaptchaStore{
		Expiration: time.Second * 180,
		PreKey:     "CAPTCHA_",
	}
}

type CaptchaStore struct {
	Expiration time.Duration
	PreKey     string
}

func (rs *CaptchaStore) Set(id string, value string) error {
	err := global.CACHE.SetCache(rs.PreKey+id, value, rs.Expiration)
	if err != nil {
		global.LOG.Error("CaptchaStore set error", zap.Error(err))
	}
	return err
}

func (rs *CaptchaStore) Get(key string, clear bool) string {
	val, exist := global.CACHE.GetCache(key)
	if clear && exist {
		global.CACHE.DeleteCache(key)
	}
	return val
}

func (rs *CaptchaStore) Verify(id, answer string, clear bool) bool {
	key := rs.PreKey + id
	v := rs.Get(key, clear)
	return v == answer
}
