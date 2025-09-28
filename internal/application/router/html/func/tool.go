package _func

import (
	"encoding/base64"
	"html/template"
	"strconv"
	"strings"
	"time"
)

func (f BaseFunc) TimeFormat(time time.Time, format string) string {
	return time.Format(format)
}

func (f BaseFunc) TimeAgo(t time.Time) string {
	unixMilli := time.Now().Unix() - t.Unix()
	switch {
	case unixMilli >= 31104000:
		return strconv.FormatInt(unixMilli/31104000, 10) + " 年前"
	case unixMilli >= 2592000:
		return strconv.FormatInt(unixMilli/2592000, 10) + " 个月前"
	case unixMilli >= 86400*2:
		return strconv.FormatInt(unixMilli/(86400*2), 10) + " 天前"
	case unixMilli >= 86400:
		return "昨天"
	case unixMilli >= 3600:
		return strconv.FormatInt(unixMilli/3600, 10) + " 小时前"
	case unixMilli >= 60:
		return strconv.FormatInt(unixMilli/60, 10) + " 分钟前"
	case unixMilli > 0:
		return strconv.FormatInt(unixMilli, 10) + " 秒前"
	default:
		return "刚刚"
	}
}

// 生成分页项
// page: 当前页，从0开始
// totalPage: 总页数
// displayCount: 每屏展示的页数
// isFull：是否展示首尾页码
func (f BaseFunc) Pagination(page int, totalPage int, displayCount int, isFull bool) []int {
	if totalPage < displayCount {
		result := make([]int, totalPage)
		for i := 0; i < totalPage; i++ {
			result[i] = i + 1
		}
		return result
	}

	isEven := displayCount%2 == 0
	left := (displayCount - 1) / 2
	right := left
	if isEven {
		right++
	}

	result := make([]int, displayCount)
	if isFull {
		if page <= left {
			for i := 0; i < displayCount-2; i++ {
				result[i] = i + 1
			}
			result[displayCount-1] = totalPage
		} else if page > totalPage-right {
			for i := 2; i < displayCount; i++ {
				result[i] = i + totalPage - displayCount + 1
			}
			result[0] = 1
		} else {
			for i := 2; i < displayCount-2; i++ {
				result[i] = i + page - left + 1
			}
			result[0] = 1
			result[displayCount-1] = totalPage
		}
	} else {
		if page <= left {
			for i := 0; i < displayCount; i++ {
				result[i] = i + 1
			}
		} else if page > totalPage-right {
			for i := 0; i < displayCount; i++ {
				result[i] = i + totalPage - displayCount + 1
			}
		} else {
			for i := 0; i < displayCount; i++ {
				result[i] = i + page - left + 1
			}
		}
	}
	return result
}

func (f BaseFunc) Html(html string) interface{} {
	return template.HTML(html)
}

func (f BaseFunc) Css(html string) interface{} {
	return template.CSS(html)
}

func (f BaseFunc) Js(html string) interface{} {
	return template.JS(html)
}

func (f BaseFunc) Base64Encode(val string) string {
	return base64.StdEncoding.EncodeToString([]byte(val))
}

func (f BaseFunc) Br(val string) string {
	return strings.ReplaceAll(val, "\n", "<br/>")
}

func (f BaseFunc) Default(value, defaultVal any) any {
	if value == nil || value == "" || value == 0.0 {
		return defaultVal
	}
	return value
}

func (f BaseFunc) Switch(flag bool, val1, val2 any) any {
	if flag {
		return val1
	}
	return val2
}

func (f BaseFunc) Params(params ...any) any {
	return params
}
