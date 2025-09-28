package enum

import "errors"

// NORMAL：正常
// FAILURE：失效
var COOKIE_STATUS = [...]string{"", "NORMAL", "FAILURE"}

type CookieStatus uint

// JsonDate反序列化
func (t *CookieStatus) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range COOKIE_STATUS {
		if status == value {
			*t = CookieStatus(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t CookieStatus) MarshalJSON() ([]byte, error) {
	return []byte("\"" + COOKIE_STATUS[t] + "\""), nil
}

func CookieStatusValue(value string) (CookieStatus, error) {
	for i, status := range COOKIE_STATUS {
		if status == value {
			return CookieStatus(i), nil
		}
	}
	return CookieStatus(0), errors.New("未找到状态码：" + value)
}
