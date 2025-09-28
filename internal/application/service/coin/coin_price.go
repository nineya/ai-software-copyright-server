package credits

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"sync"
)

type CreditsPriceService struct {
	service.BaseService
}

var onceBuy = sync.Once{}
var creditsPriceService *CreditsPriceService

// 获取单例
func GetCreditsPriceService() *CreditsPriceService {
	onceBuy.Do(func() {
		creditsPriceService = new(CreditsPriceService)
		creditsPriceService.Db = global.DB
	})
	return creditsPriceService
}

func (s *CreditsPriceService) GetAll() ([]table.CreditsPrice, error) {
	list := make([]table.CreditsPrice, 0)
	err := s.Db.Find(&list)
	return list, err
}
