package enum

import (
	"github.com/pkg/errors"
)

var LOG_TYPE = [...]string{
	"",                     // 0, 未定义
	"SITE_INITIALIZED",     // 1,站点初始化
	"SITE_STATUS_UPDATE",   // 2,站点状态变更
	"STATIC_STORE",         // 3,静态存储
	"",                     // 4
	"",                     // 5
	"",                     // 6
	"",                     // 7
	"LOG_CLEAR",            // 8,日志清空
	"CHANGE_LICENSE",       // 9，更换 LICENSE
	"",                     // 10
	"",                     // 11
	"DENIAL_ACCESS",        // 12,拒绝访问
	"",                     // 13
	"",                     // 14
	"",                     // 15
	"",                     // 16
	"",                     // 17
	"",                     // 18
	"",                     // 19
	"",                     // 20
	"",                     // 21
	"ADMIN_LOGIN",          // 22,管理员登录
	"ADMIN_LOGOUT",         // 23,管理员注销登录
	"ADMIN_PROFILE",        // 24,管理员资料编辑
	"",                     // 25
	"",                     // 26
	"",                     // 27
	"",                     // 28
	"",                     // 29
	"FAILED_LOGIN",         // 30,登录失败
	"",                     // 31
	"",                     // 32
	"",                     // 33
	"",                     // 34
	"",                     // 35
	"",                     // 36
	"",                     // 37
	"",                     // 38
	"",                     // 39
	"CREATED_ADMIN",        // 40,创建管理员
	"CREATED_ADMIN_GROUP",  // 41,创建管理员分组
	"CREATED_POST",         // 42,创建新文章
	"CREATED_CATEGORY",     // 43,创建文章分类
	"CREATED_SHEET",        // 44,创建自建页面
	"CREATED_MENU",         // 45,创建菜单
	"CREATED_TAG",          // 46,创建文章标签
	"CREATED_FROM_SCHEMA",  // 47,创建收集表单
	"CREATED_LINK",         // 48,创建友链
	"CREATED_OPTION",       // 49,创建系统变量
	"CREATED_PHOTO",        // 50,创建图库图片
	"",                     // 51
	"",                     // 52
	"UPLOAD_ATTACHMENT",    // 53,上传附件
	"UPLOAD_THEME",         // 54,上传主题
	"",                     // 55
	"",                     // 56
	"DELETED_ADMIN",        // 57,删除管理员
	"DELETED_ADMIN_GROUP",  // 58,删除管理员分组
	"DELETED_POST",         // 59,删除文章
	"DELETED_ATTACHMENT",   // 60,删除附件
	"DELETED_CATEGORY",     // 61,删除文章分类
	"DELETED_SHEET",        // 62,删除自建页面
	"DELETED_MENU",         // 63,删除菜单
	"DELETED_TAG",          // 64,删除文章标签
	"DELETED_THEME",        // 65,删除主题
	"DELETED_FROM_SCHEMA",  // 66,删除收集表单
	"DELETED_LINK",         // 67,删除友链
	"DELETED_OPTION",       // 68,删除系统变量
	"DELETED_PHOTO",        // 69,删除图库图片
	"",                     // 70
	"",                     // 71
	"UPDATED_ADMIN",        // 72,更新管理员
	"UPDATED_ADMIN_GROUP",  // 73,更新管理员分组
	"UPDATED_ATTACHMENT",   // 74,更新附件信息
	"UPDATED_CATEGORY",     // 75,更新文章分类
	"UPDATED_MENU",         // 76,更新菜单
	"UPDATED_TAG",          // 77,更新文章标签
	"UPDATED_SITE",         // 78,更新站点信息
	"UPDATED_THEME",        // 79,更新主题
	"UPDATED_FROM_SCHEMA",  // 80,更新收集表单
	"UPDATED_LINK",         // 81,更新友链
	"UPDATED_OPTION",       // 82,更新系统变量
	"UPDATED_PHOTO",        // 83,更新图库图片
	"",                     // 84
	"",                     // 85
	"EDITED_POST",          // 86,修改文章
	"EDITED_SHEET",         // 87,修改自建页面
	"EDITED_THEME_SETTING", // 88,修改主题配置
	"EDITED_THEME",         // 89，修改主题文件
}

type LogType uint

// JsonDate反序列化
func (t *LogType) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range LOG_TYPE {
		if status == value {
			*t = LogType(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t LogType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + LOG_TYPE[t] + "\""), nil
}

func LogTypeValue(value string) (LogType, error) {
	for i, postType := range LOG_TYPE {
		if postType == value {
			return LogType(i), nil
		}
	}
	return LogType(0), errors.New("未找到状态码：" + value)
}
