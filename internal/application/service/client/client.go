package client

import (
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"sync"
)

type ClientService struct {
	service.BaseService
	awaitResultMap map[string]chan any
}

var onceClient = sync.Once{}
var clientService *ClientService

// 获取单例
func GetClientService() *ClientService {
	onceClient.Do(func() {
		clientService = new(ClientService)
		clientService.Db = global.DB
	})
	return clientService
}
