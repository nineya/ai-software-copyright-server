package response

import (
	"time"
	"tool-server/internal/application/model/table"
)

type GeneralStatisticResponse struct {
	PostCount    int64     `json:"postCount"`
	CommentCount int64     `json:"commentCount"`
	VisitCount   int64     `json:"visitCount"`
	ExpireDate   time.Time `json:"expireDate"`
	ExpireDays   int       `json:"expireDays"`
}

type GeneralStatisticWithAdminResponse struct {
	GeneralStatisticResponse
	Admin table.Admin `json:"admin"`
}

type ChartStatisticResponse struct {
	Date  string `json:"date"`
	Label string `json:"label"`
	Value int    `json:"value"`
}
