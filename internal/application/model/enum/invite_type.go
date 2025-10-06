package enum

import (
	"github.com/pkg/errors"
)

var INVITE_TYPE = [...]string{
	"",         // 0, 未定义
	"REGISTER", // 1,注册
	"ACTIVE",   // 2,活跃
}

type InviteType uint

// JsonDate反序列化
func (t *InviteType) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range INVITE_TYPE {
		if status == value {
			*t = InviteType(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t InviteType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + INVITE_TYPE[t] + "\""), nil
}

func InviteTypeValue(value string) (InviteType, error) {
	for i, postType := range INVITE_TYPE {
		if postType == value {
			return InviteType(i), nil
		}
	}
	return InviteType(0), errors.New("未找到状态码：" + value)
}
