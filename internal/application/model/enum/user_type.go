package enum

import (
	"github.com/pkg/errors"
)

var USER_TYPE = [...]string{
	"",      // 0,未定义
	"ADMIN", // 1,用户
	"USER",  // 2,管理员
}

type UserType uint

// JsonDate反序列化
func (t *UserType) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range USER_TYPE {
		if status == value {
			*t = UserType(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t UserType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + USER_TYPE[t] + "\""), nil
}

func UserTypeValue(value string) (UserType, error) {
	for i, postType := range USER_TYPE {
		if postType == value {
			return UserType(i), nil
		}
	}
	return UserType(0), errors.New("未找到状态码：" + value)
}
