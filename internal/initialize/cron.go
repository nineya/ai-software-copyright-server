package initialize

import (
	"github.com/robfig/cron/v3"
	"tool-server/internal/global"
	"tool-server/internal/task"
	"tool-server/internal/utils"
)

func InitCron() {
	cronTab := cron.New()
	// 每天2点执行清理请求统计数据的定时任务
	_, err := cronTab.AddFunc("0 2 * * ?", task.ClearTask)
	utils.PanicErr(err)

	cronTab.Start()
	global.CRON = cronTab
}
