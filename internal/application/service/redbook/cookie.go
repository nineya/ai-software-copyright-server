package redbook

import (
	"sync"
	"tool-server/internal/application/model/table"
	"tool-server/internal/application/service"
	"tool-server/internal/global"
)

type CookieService struct {
	service.CrudService[table.RedbookCookie]
}

var onceUser = sync.Once{}
var cookieService *CookieService

// 获取单例
func GetCookieService() *CookieService {
	onceUser.Do(func() {
		cookieService = new(CookieService)
		cookieService.Db = global.DB
	})
	return cookieService
}
