package main

import (
	"tool-server/internal/global"
	"tool-server/internal/initialize"
)

//go:generate swag init -g index.go -d internal/application/router/api/admin,internal/application/param,internal/application/model -o ./docs/admin --instanceName=admin
//go:generate swag init -g index.go -d internal/application/router/api/content,internal/application/param,internal/application/model -o ./docs/content --instanceName=content
func main() {
	initialize.InitSystemConfig()
	initialize.InitLogger()
	initialize.InitDatabase()
	initialize.InitRedis()
	initialize.InitCache()
	initialize.InitCron()
	router := initialize.InitRouter()
	initialize.RunServer(router)
	defer Close()
}

func Close() {
	if global.DB != nil {
		// 程序结束前关闭数据库链接
		global.DB.Close()
		global.LOG.Info("Close database connection.")
	}
	if global.LOG != nil {
		global.LOG.Sync()
	}
	if global.CRON != nil {
		global.CRON.Stop()
	}
}
