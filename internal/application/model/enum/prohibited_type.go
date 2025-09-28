package enum

import (
	"github.com/pkg/errors"
)

var PROHIBITED_TYPE = [...]string{
	"",           // 0, 未定义
	"PROHIBITED", // 1,违禁
	"SENSITIVE",  // 2,敏感
	"CUSTOM",     // 3,自定义
}

type ProhibitedType uint

// JsonDate反序列化
func (t *ProhibitedType) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range PROHIBITED_TYPE {
		if status == value {
			*t = ProhibitedType(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t ProhibitedType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + PROHIBITED_TYPE[t] + "\""), nil
}

func ProhibitedTypeValue(value string) (ProhibitedType, error) {
	for i, postType := range PROHIBITED_TYPE {
		if postType == value {
			return ProhibitedType(i), nil
		}
	}
	return ProhibitedType(0), errors.New("未找到状态码：" + value)
}
