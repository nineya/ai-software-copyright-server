package utils

import (
	"crypto/rand"
	"encoding/base64"
	"math"
	"net/url"
	"regexp"
	"strings"
)

func GetHost(rawUrl string) string {
	if rawUrl == "" {
		return ""
	}
	newUrl, _ := url.Parse(rawUrl)
	return newUrl.Host
}

func RemoveSuffix(str string, suffix string) string {
	if strings.HasSuffix(str, suffix) {
		return str[:len(str)-len(suffix)]
	} else {
		return str // 不需要修改直接返回原始字符串
	}
}

func MaskContent(content string) string {
	contentRune := []rune(content)
	switch len(contentRune) {
	case 1:
		return "*"
	case 2:
		return string(contentRune[0:1]) + "*"
	default:
		maskSize := int(math.Round(float64(len(contentRune)) / 3))
		var str string
		for i := 0; i < maskSize; i++ {
			str += "*"
		}
		return string(contentRune[0:maskSize]) + str + string(contentRune[maskSize+maskSize:])
	}
}

// 判断是否为手机号的函数
func CheckPhone(phone string) bool {
	phoneRegexp := regexp.MustCompile(`^1[0-9]{10,13}$`)
	return phoneRegexp.MatchString(phone)
}

// 判断是否为邮箱的函数
func CheckEmail(email string) bool {
	emailRegexp := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegexp.MatchString(email)
}

// 生成随机字符串
func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:length], nil
}
