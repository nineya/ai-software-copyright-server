package task

import (
	"ai-software-copyright-server/internal/application/model/table"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/global"
	"fmt"
	"time"
)

// 过期时间通知
func ExpireNoticeTask() {
	global.LOG.Info("定时任务（服务过期通知）：开始执行")
	// 短链服务过期提醒
	slConfigures := make([]table.NetdiskShortLinkConfigure, 0)
	err := global.DB.Where("DATE(custom_expire_time) = CURDATE() + INTERVAL 15 DAY or DATE(custom_expire_time) = CURDATE() + INTERVAL 3 DAY or DATE(custom_expire_time) = CURDATE()").Find(&slConfigures)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("定时任务（服务过期通知）：获取网盘过期短链列表失败: %+v", err))
	}
	for _, c := range slConfigures {
		sendNotice(c.UserId, c.CustomExpireTime, "网盘短链定制版", "失效后网盘短链自定义域名将无法使用")
	}
	// 搜索网站过期提醒
	ssConfigures := make([]table.NetdiskSearchSiteConfigure, 0)
	err = global.DB.Where("DATE(expire_time) = CURDATE() + INTERVAL 15 DAY or DATE(expire_time) = CURDATE() + INTERVAL 3 DAY or DATE(expire_time) = CURDATE()").Find(&ssConfigures)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("定时任务（服务过期通知）：获取过期网盘搜索网站列表失败: %+v", err))
	}
	for _, c := range ssConfigures {
		sendNotice(c.UserId, c.ExpireTime, "网盘搜索网站云服务版", "失效后网盘搜索网站将无法访问")
	}
	// 网盘助手过期提醒
	hConfigures := make([]table.NetdiskHelperConfigure, 0)
	err = global.DB.Where("DATE(expire_time) = CURDATE() + INTERVAL 15 DAY or DATE(expire_time) = CURDATE() + INTERVAL 3 DAY or DATE(expire_time) = CURDATE()" +
		"or DATE(wechat_expire_time) = CURDATE() + INTERVAL 15 DAY or DATE(wechat_expire_time) = CURDATE() + INTERVAL 3 DAY or DATE(wechat_expire_time) = CURDATE()").Find(&hConfigures)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("定时任务（服务过期通知）：获取过期网盘助手列表失败: %+v", err))
	}
	for _, c := range hConfigures {
		if c.WechatExpireTime != nil {
			expireDay := getDaysBetweenUnix(time.Now(), *c.WechatExpireTime)
			if expireDay > 0 && expireDay <= 15 {
				sendNotice(c.UserId, c.WechatExpireTime, "网盘助手微信工具人版", "失效后微信群资料搜索机器人、网盘资源采集和网盘管理等功能将无法使用")
				continue
			}
		}
		if c.ExpireTime != nil {
			expireDay := getDaysBetweenUnix(time.Now(), *c.ExpireTime)
			if expireDay > 0 && expireDay <= 15 {
				sendNotice(c.UserId, c.ExpireTime, "网盘助手", "失效后网盘资源采集、网盘管理等功能将无法使用")
				continue
			}
		}
	}
	global.LOG.Info("定时任务（服务过期通知）：执行完成")
}

// 发送通知
func sendNotice(userId int64, expireTime *time.Time, serviceName, serviceDesc string) {
	if expireTime == nil {
		return
	}
	expireDay := getDaysBetweenUnix(time.Now(), *expireTime)
	var expireText string
	if expireDay > 0 {
		expireText = fmt.Sprintf("将在%d天后过期", expireDay)
	} else {
		expireText = "即将过期"
	}
	title := fmt.Sprintf("服务过期通知：您订阅的%s服务%s", serviceName, expireText)
	content := fmt.Sprintf("您好：\n    您订阅的%s服务%s，%s。如需再次订阅，联系VX（玖涯菜菜子）：nineyaccz", serviceName, expireText, serviceDesc)
	err := userSev.GetUserService().SendMail(userId, title, content)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("定时任务（服务过期通知）：发送%s过期通知失败: %+v", serviceName, err))
	}
}

// 取得两个日期的时间间隔
func getDaysBetweenUnix(t1, t2 time.Time) int {
	// 转换为当天的起始时间
	startOfDay1 := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, t1.Location())
	startOfDay2 := time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, t2.Location())

	// 使用 Unix 时间戳计算天数
	seconds := startOfDay2.Unix() - startOfDay1.Unix()
	days := int(seconds / 86400) // 86400 = 24 * 60 * 60
	return days
}
