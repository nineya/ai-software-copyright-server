package user

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"sync"
)

type BuyService struct {
	service.UserCrudService[table.Buy]
}

var onceBuy = sync.Once{}
var buyService *BuyService

// 获取单例
func GetBuyService() *BuyService {
	onceBuy.Do(func() {
		buyService = new(BuyService)
		buyService.Db = global.DB
	})
	return buyService
}
