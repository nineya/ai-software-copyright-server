package main

import (
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/initialize"
	"embed"
)

//go:embed resource/*
var f embed.FS

//go:generate swag init -g index.go -d internal/application/router/api/content,internal/application/param,internal/application/model -o ./docs/content --instanceName=content
func main() {
	global.FS = f
	initialize.InitSystemConfig()
	initialize.InitLogger()
	initialize.InitDatabase()
	initialize.InitRedis()
	initialize.InitCache()
	initialize.InitWechatPay()
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
