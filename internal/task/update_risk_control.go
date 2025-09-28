package task

import (
	rcSev "ai-software-copyright-server/internal/application/service/risk_control"
	"ai-software-copyright-server/internal/global"
)

func UpdateRiskControlTask() {
	global.LOG.Info("更新风控数据：开始执行定时任务")
	err := rcSev.GetRiskControlService().UpdateRiskControlScore()
	if err != nil {
		global.LOG.Sugar().Errorf("更新风控数据：执行定时任务失败，%+v", err)
		return
	}
	global.LOG.Info("更新风控数据：定时任务执行完成")
}
