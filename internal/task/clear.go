package task

import (
	statisSev "ai-software-copyright-server/internal/application/service/statistic"
	"ai-software-copyright-server/internal/global"
	"fmt"
)

func ClearTask() {
	global.LOG.Info("Clear statistic.")
	count, err := statisSev.GetStatisticService().ClearStatistic()
	if err != nil {
		global.LOG.Sugar().Error("Clear statistic failure.", err)
		return
	}
	global.LOG.Info(fmt.Sprintf("Clear %d Statistic data.", count))
}
