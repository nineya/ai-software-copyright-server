package task

import (
	"fmt"
	statisSev "tool-server/internal/application/service/statistic"
	"tool-server/internal/global"
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
