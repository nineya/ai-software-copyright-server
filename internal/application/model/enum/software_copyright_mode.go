package enum

import (
	"github.com/pkg/errors"
)

var SOFTWARE_COPYRIGHT_MODE = [...]string{
	"",         // 0,未定义
	"ALL",      // 1,全部
	"REQUEST",  // 2,申请表
	"CODE",     // 3,代码材料
	"DOCUMENT", // 4,文档材料
}

type SoftwareCopyrightMode uint

// JsonDate反序列化
func (t *SoftwareCopyrightMode) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range SOFTWARE_COPYRIGHT_MODE {
		if status == value {
			*t = SoftwareCopyrightMode(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t SoftwareCopyrightMode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + SOFTWARE_COPYRIGHT_MODE[t] + "\""), nil
}

func SoftwareCopyrightModeValue(value string) (SoftwareCopyrightMode, error) {
	for i, postType := range SOFTWARE_COPYRIGHT_MODE {
		if postType == value {
			return SoftwareCopyrightMode(i), nil
		}
	}
	return SoftwareCopyrightMode(0), errors.New("未找到状态码：" + value)
}
