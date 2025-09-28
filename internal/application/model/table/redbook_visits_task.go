package table

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"time"
)

type RedbookVisitsTask struct {
	Id           int64           `json:"id" xorm:"<- PK AUTOINCR comment('任务id')" label:"任务id"`
	AdminId      int64           `json:"adminId,omitempty" xorm:"notnull comment('管理员id')" label:"管理员id"`
	Name         string          `json:"name" xorm:"VARCHAR(50) notnull comment('任务名称')" binding:"lte=50" label:"任务名称"`
	NoteId       string          `json:"noteId" xorm:"VARCHAR(50) notnull comment('小红书文章id')" label:"小红书文章id"`
	AuthorId     string          `json:"authorId" xorm:"VARCHAR(50) notnull comment('文章作者id')" label:"文章作者id"`
	CurrentCount int             `json:"currentCount" xorm:"INT notnull comment('当前计数')" label:"当前计数"`
	TargetCount  int             `json:"targetCount" xorm:"INT notnull comment('目标计数')" label:"目标计数"`
	IntervalTime int             `json:"intervalTime" xorm:"INT notnull comment('间隔时长(秒)')" label:"间隔时长(秒)"`
	Status       enum.TaskStatus `json:"status" xorm:"SMALLINT notnull comment('任务状态')" label:"任务状态"`
	CreateTime   *time.Time      `json:"createTime" xorm:"DATETIME created"`
	UpdateTime   *time.Time      `json:"updateTime" xorm:"DATETIME updated"`
}

func (RedbookVisitsTask) TableName() string {
	return "redbook_visits_task"
}

func (t *RedbookVisitsTask) SetAdminId(adminId int64) {
	t.AdminId = adminId
}
