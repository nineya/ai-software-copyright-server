package enum

import (
	"github.com/pkg/errors"
)

var HINT_LEVEL = [...]string{
	"",      // 0, 未定义
	"INFO",  // 1,信息
	"NORM",  // 2,正常
	"WARN",  // 3,提醒
	"RISK",  // 4,危险
	"REFER", // 5,参考
}

type HintLevel uint

// JsonDate反序列化
func (t *HintLevel) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range HINT_LEVEL {
		if status == value {
			*t = HintLevel(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t HintLevel) MarshalJSON() ([]byte, error) {
	return []byte("\"" + HINT_LEVEL[t] + "\""), nil
}

func HintLevelValue(value string) (HintLevel, error) {
	for i, postType := range HINT_LEVEL {
		if postType == value {
			return HintLevel(i), nil
		}
	}
	return HintLevel(0), errors.New("未找到状态码：" + value)
}
