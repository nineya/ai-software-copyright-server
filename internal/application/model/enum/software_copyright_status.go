package enum

import (
	"github.com/pkg/errors"
)

var SOFTWARE_COPYRIGHT_STATUS = [...]string{
	"",         // 0,未定义
	"GENERATE", // 1,生成中
	"COMPLETE", // 2,已完成
	"FAILURE",  // 3,失败
}

type SoftwareCopyrightStatus uint

// JsonDate反序列化
func (t *SoftwareCopyrightStatus) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range SOFTWARE_COPYRIGHT_STATUS {
		if status == value {
			*t = SoftwareCopyrightStatus(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t SoftwareCopyrightStatus) MarshalJSON() ([]byte, error) {
	return []byte("\"" + SOFTWARE_COPYRIGHT_STATUS[t] + "\""), nil
}

func SoftwareCopyrightStatusValue(value string) (SoftwareCopyrightStatus, error) {
	for i, postType := range SOFTWARE_COPYRIGHT_STATUS {
		if postType == value {
			return SoftwareCopyrightStatus(i), nil
		}
	}
	return SoftwareCopyrightStatus(0), errors.New("未找到状态码：" + value)
}
