package task

import (
	netdSev "ai-software-copyright-server/internal/application/service/netdisk"
	"ai-software-copyright-server/internal/global"
	"fmt"
	"time"
)

func CheckQuarkResourceTask() {
	defer func() {
		if r := recover(); r != nil {
			global.LOG.Error(fmt.Sprintf("夸克网盘资源检测任务执行失败: %+v", r))
			time.Sleep(10 * time.Minute)
			CheckQuarkResourceTask()
		}
	}()
	global.LOG.Info("Check Quark Resource.")
	for true {
		count, err := netdSev.GetResourceService().CheckQuarkResource()
		if err != nil {
			global.LOG.Error(fmt.Sprintf("查询夸克网盘资源失败: %+v", err))
		}
		if count == 0 {
			time.Sleep(1 * time.Hour)
			global.LOG.Info("查询到待检测夸克资源为空，等待一小时")
		}
	}
}

func CheckBaiduResourceTask() {
	defer func() {
		if r := recover(); r != nil {
			global.LOG.Error(fmt.Sprintf("百度网盘资源检测任务执行失败: %+v", r))
			time.Sleep(10 * time.Minute)
			CheckQuarkResourceTask()
		}
	}()
	global.LOG.Info("Check Baidu Resource.")
	for true {
		count, err := netdSev.GetResourceService().CheckBaiduResource()
		if err != nil {
			global.LOG.Error(fmt.Sprintf("查询百度网盘资源失败: %+v", err))
		}
		if count == 0 {
			time.Sleep(1 * time.Hour)
			global.LOG.Info("查询到待检测百度网盘资源为空，等待一小时")
		}
	}
}
