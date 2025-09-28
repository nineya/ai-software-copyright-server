package initialize

import (
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/task"
	"ai-software-copyright-server/internal/utils"
	"github.com/robfig/cron/v3"
)

func InitCron() {
	cronTab := cron.New()
	// 每天2点执行清理请求统计数据的定时任务
	_, err := cronTab.AddFunc("0 2 * * ?", task.ClearTask)
	utils.PanicErr(err)
	//// 每天1点执行风控数据更新 TODO 去掉风控
	//_, err = cronTab.AddFunc("0 1 * * ?", task.UpdateRiskControlTask)
	//utils.PanicErr(err)
	// 每天10点执行服务过期通知
	_, err = cronTab.AddFunc("0 10 * * ?", task.ExpireNoticeTask)
	utils.PanicErr(err)
	// 每隔5分钟执行Socket心跳
	_, err = cronTab.AddFunc("*/5 * * * ?", task.SocketHeartbeatTask)
	utils.PanicErr(err)
	// 循环执行夸克网盘资源分析
	go task.CheckQuarkResourceTask()
	// 循环执行百度网盘资源分析
	go task.CheckBaiduResourceTask()

	cronTab.Start()
	global.CRON = cronTab
}
