package response

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"time"
)

type TimeClockCreateResponse struct {
	UserBuyResponse
	Data table.TimeClock `json:"data"`
}

type TimeClockMyInfoResponse struct {
	table.TimeClock
	Member       table.TimeClockMember `json:"member"`
	ClockIn      bool                  `json:"clockIn"`
	ClockInCount int64                 `json:"clockInCount"`
}

type TimeClockMemberInfoResponse struct {
	Id           int64                      `json:"id,omitempty"` //主键
	UserId       int64                      `json:"userId,omitempty"`
	ClockId      int64                      `json:"clockId"` //主键
	Status       enum.TimeClockMemberStatus `json:"status"`
	JoinTime     *time.Time                 `json:"joinTime"`
	ClockInCount int                        `json:"clockInCount"`
	Avatar       string                     `json:"avatar,omitempty"`
	Nickname     string                     `json:"nickname,omitempty"` //昵称
	InviteCode   string                     `json:"inviteCode,omitempty"`
	Remark       string                     `json:"remark,omitempty"`
}
