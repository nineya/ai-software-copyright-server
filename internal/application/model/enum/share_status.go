package enum

import "errors"

var SHARE_STATUS = [...]string{
	"",
	"AWAIT", // 等待
	"PASS",  // 通过
	"DENY",  // 拒绝
}

type ShareStatus uint

// JsonDate反序列化
func (t *ShareStatus) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range SHARE_STATUS {
		if status == value {
			*t = ShareStatus(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t ShareStatus) MarshalJSON() ([]byte, error) {
	return []byte("\"" + SHARE_STATUS[t] + "\""), nil
}

func ShareStatusValue(value string) (ShareStatus, error) {
	for i, status := range SHARE_STATUS {
		if status == value {
			return ShareStatus(i), nil
		}
	}
	return ShareStatus(0), errors.New("未找到状态码：" + value)
}
