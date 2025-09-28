package response

import "ai-software-copyright-server/internal/application/model/table"

type ShortLinkRedirectResponse struct {
	UserBuyResponse
	Url string `json:"url"`
}

type ShortLinkStatisticResponse struct {
	UserBuyResponse
	Today        table.ShortLinkStatistic   `json:"today"`        // 今天数据
	TodayOrigins []table.ShortLinkStatistic `json:"todayOrigins"` // 今天访问来源
	Total        table.ShortLinkStatistic   `json:"total"`
	TotalOrigins []table.ShortLinkStatistic `json:"totalOrigins"` // 总访问来源
	Days         []table.ShortLinkStatistic `json:"days"`
}
