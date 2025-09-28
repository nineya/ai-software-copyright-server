package enum

import "errors"

// NORMAL：正常
// FAILURE：失效
var TIME_CLOCK_MEMBER_STATUS = [...]string{
	"",       // 0,未定义
	"NORMAL", // 1,正常
	"AUDIT",  // 2,待审核
}

type TimeClockMemberStatus uint

// JsonDate反序列化
func (t *TimeClockMemberStatus) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range TIME_CLOCK_MEMBER_STATUS {
		if status == value {
			*t = TimeClockMemberStatus(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t TimeClockMemberStatus) MarshalJSON() ([]byte, error) {
	return []byte("\"" + TIME_CLOCK_MEMBER_STATUS[t] + "\""), nil
}

func TimeClockMemberStatusValue(value string) (TimeClockMemberStatus, error) {
	for i, status := range TIME_CLOCK_MEMBER_STATUS {
		if status == value {
			return TimeClockMemberStatus(i), nil
		}
	}
	return TimeClockMemberStatus(0), errors.New("未找到状态码：" + value)
}
