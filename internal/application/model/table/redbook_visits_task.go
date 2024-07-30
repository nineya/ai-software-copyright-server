package table

import "time"

type RedbookVisitsTask struct {
	Id           int64      `json:"id" xorm:"<- PK AUTOINCR comment('任务id')"`
	AdminId      int64      `json:"adminId,omitempty" xorm:"-> notnull comment('管理员id')"`
	Name         string     `json:"name" xorm:"VARCHAR(50) notnull comment('任务名称')" binding:"lte=50"`
	NoteId       int64      `json:"noteId" xorm:"-> VARCHAR(50) notnull comment('小红书文章id')"`
	CurrentCount int        `json:"currentCount" xorm:"INT notnull comment('当前计数')"`
	TargetCount  int        `json:"targetCount" xorm:"INT notnull comment('目标计数')"`
	Status       int        `json:"status" xorm:"SMALLINT notnull comment('任务状态')"`
	CreateTime   *time.Time `json:"createTime" xorm:"DATETIME created"`
	UpdateTime   *time.Time `json:"updateTime" xorm:"DATETIME updated"`
}

func (RedbookVisitsTask) TableName() string {
	return "redbook_visits_task"
}
