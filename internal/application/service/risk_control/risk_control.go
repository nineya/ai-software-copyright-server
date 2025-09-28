package risk_control

import (
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"sync"
)

type RiskControlService struct {
	service.BaseService
}

var onceUser = sync.Once{}
var riskControlService *RiskControlService

// 获取单例
func GetRiskControlService() *RiskControlService {
	onceUser.Do(func() {
		riskControlService = new(RiskControlService)
		riskControlService.Db = global.DB
	})
	return riskControlService
}

// 更新风控得分
func (s *RiskControlService) UpdateRiskControlScore() error {
	// 100分是基础分
	// 对所有用户，每日自然扣减风控分
	_, err := s.Db.Exec(`update user set risk_control_score = risk_control_score - 5 where risk_control_score > 100`)
	if err != nil {
		global.LOG.Sugar().Warnf("扣减分值超过100用户的风控得分失败: %+v", err)
	}
	_, err = s.Db.Exec(`update user set risk_control_score = risk_control_score - 1 where risk_control_score <= 100`)
	if err != nil {
		global.LOG.Sugar().Warnf("扣减分值小于等于100用户的风控得分失败: %+v", err)
	}
	// 短链访问增加风控得分，近一个月的数据之和除100
	_, err = s.Db.Exec(`UPDATE user
INNER JOIN (select user_id, sum(visits)/100 as num from short_link where create_time > DATE_SUB(CURDATE(), INTERVAL 1 MONTH) GROUP BY user_id) t
ON user.id = t.user_id
SET user.risk_control_score = user.risk_control_score + num
where risk_control_score > -99999999`)
	if err != nil {
		global.LOG.Sugar().Warnf("通过短链访问更新风控得分失败: %+v", err)
	}
	// 根据昨日的币变动更新风控得分
	_, err = s.Db.Exec(`UPDATE user
INNER JOIN (select user_id,sum(case type 
when 1 then change_credits / -2
when 2 then change_credits
when 3 then change_credits * 2
when 4 then change_credits * 10
when 5 then change_credits / 2
end) num from credits_change where  DATE(create_time) = DATE_SUB(CURDATE(), INTERVAL 1 DAY) GROUP BY user_id) t
ON user.id = t.user_id
SET user.risk_control_score = user.risk_control_score + num
where risk_control_score > -99999999`)
	if err != nil {
		global.LOG.Sugar().Warnf("通过昨日的积分变动更新风控得分失败: %+v", err)
	}
	// 通过近5天的接口访问次数更新风控得分
	_, err = s.Db.Exec(`UPDATE user
INNER JOIN (select user_id,if(count(*) > 500, 500-count(*),count(*))/10 num from user_log where create_time > DATE_SUB(CURDATE(), INTERVAL 5 DAY) GROUP BY user_id) t
ON user.id = t.user_id
SET user.risk_control_score = user.risk_control_score + num
where risk_control_score > -99999999`)
	if err != nil {
		global.LOG.Sugar().Warnf("通过近5天的接口访问次数更新风控得分失败: %+v", err)
	}
	return err
}
