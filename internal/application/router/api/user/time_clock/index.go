package time_clock

type RouterGroup struct {
	TimeClockApiRouter
	TimeClockMemberApiRouter
	TimeClockRecordApiRouter
}
