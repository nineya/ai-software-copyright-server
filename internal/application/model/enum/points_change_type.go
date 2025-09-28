package enum

import (
	"github.com/pkg/errors"
)

var POINTS_CHANGE_TYPE = [...]string{
	"",        // 0,未定义
	"BUY",     // 1,购买
	"REWARD",  // 2,激励
	"INVITER", // 3,邀请
	"PAY",     // 4,充值
	"GIVE",    // 5,送币
}

type CreditsChangeType uint

// JsonDate反序列化
func (t *CreditsChangeType) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range POINTS_CHANGE_TYPE {
		if status == value {
			*t = CreditsChangeType(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t CreditsChangeType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + POINTS_CHANGE_TYPE[t] + "\""), nil
}

func CreditsChangeTypeValue(value string) (CreditsChangeType, error) {
	for i, postType := range POINTS_CHANGE_TYPE {
		if postType == value {
			return CreditsChangeType(i), nil
		}
	}
	return CreditsChangeType(0), errors.New("未找到状态码：" + value)
}
