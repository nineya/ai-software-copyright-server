package enum

import (
	"github.com/pkg/errors"
)

var STUDY_TYPE = [...]string{
	"",        // 0, 未定义
	"REDBOOK", // 1,小红书
	"MP",      // 2,公众号
}

type StudyType uint

// JsonDate反序列化
func (t *StudyType) UnmarshalJSON(data []byte) error {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range STUDY_TYPE {
		if status == value {
			*t = StudyType(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t StudyType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + STUDY_TYPE[t] + "\""), nil
}

func StudyTypeValue(value string) (StudyType, error) {
	for i, postType := range STUDY_TYPE {
		if postType == value {
			return StudyType(i), nil
		}
	}
	return StudyType(0), errors.New("未找到状态码：" + value)
}
