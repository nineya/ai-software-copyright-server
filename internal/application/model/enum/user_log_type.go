package enum

import (
	"github.com/pkg/errors"
)

var USER_LOG_TYPE = [...]string{
	"",                             // 0, 未定义
	"SITE_INITIALIZED",             // 1,站点初始化
	"SITE_STATUS_UPDATE",           // 2,站点状态变更
	"STATIC_STORE",                 // 3,静态存储
	"",                             // 4
	"",                             // 5
	"",                             // 6
	"",                             // 7
	"LOG_CLEAR",                    // 8,日志清空
	"CHANGE_LICENSE",               // 9，更换 LICENSE
	"",                             // 10
	"",                             // 11
	"DENIAL_ACCESS",                // 12,拒绝访问
	"",                             // 13
	"",                             // 14
	"",                             // 15
	"",                             // 16
	"",                             // 17
	"",                             // 18
	"",                             // 19
	"",                             // 20
	"",                             // 21
	"USER_LOGIN",                   // 22,用户登录
	"USER_LOGOUT",                  // 23,用户注销登录
	"USER_PROFILE",                 // 24,用户资料编辑
	"USER_ACCESS_KEY_UPDATE",       // 25，用户AccessKey更新
	"",                             // 26
	"",                             // 27
	"",                             // 28
	"",                             // 29
	"FAILED_LOGIN",                 // 30,登录失败
	"",                             // 31
	"",                             // 32
	"",                             // 33
	"",                             // 34
	"",                             // 35
	"",                             // 36
	"POINTS_ORDER_CREATE",          // 37,创建付款订单
	"MP_IMAGE_TEXT_OPTIMIZE",       // 38,公众号图文优化
	"NETDISK_CONFIGURE_SAVE",       // 39,保存网盘资源配置
	"NETDISK_RESOURCE_CREATE",      // 40,创建网盘资源
	"NETDISK_RESOURCE_IMPORT",      // 41,批量导入网盘资源
	"NETDISK_RESOURCE_DELETE",      // 42,删除网盘资源
	"NETDISK_RESOURCE_UPDATE",      // 43,更新网盘资源
	"NETDISK_SHORT_LINK_CREATE",    // 44,创建网盘短链
	"NETDISK_SHORT_LINK_REDIRECT",  // 45,网盘短链重定向
	"WX_PAY_NOTIFY",                // 46,微信支付回调
	"QRCODE_BUILD",                 // 47,创建收集表单
	"REDBOOK_PROHIBITED_DETECTION", // 48,创建友链
	"REDBOOK_PROHIBITED_CUSTOM",    // 49,小红书用户自定义敏感词
	"REDBOOK_REMOVE_WATERMARK",     // 50,小红书去水印
	"REDBOOK_VALUATION",            // 51,小红书账号估值
	"REDBOOK_WEIGHT",               // 52,小红书账号权重分析
	"REDBOOK_WRITE_TITLE",          // 53,小红书爆款标题生成
	"REDBOOK_WRITE_NOTE",           // 54,小红书笔记帮写/优化
	"REDBOOK_WRITE_PLANTING",       // 55,小红书种草笔记生成
	"USER_INFO_UPDATE",             // 56,用户信息更新
	"DELETED_NETDISK_RESOURCE",     // 57,删除网盘资源
	"NETDISK_SHORT_LINK_STATISTIC", // 58,网盘短链分析
	"NETDISK_HELPER_SEND_REQUEST",  // 59,发送网盘助手请求
	"QRCODE_LOOSE_CREATE",          // 60,创建活码
	"QRCODE_LOOSE_DELETE",          // 61,删除活码
	"QRCODE_LOOSE_UPDATE",          // 62,更新活码
	"QRCODE_LOOSE_ADD_IMAGE",       // 63,活码添加图片
	"QRCODE_LOOSE_DELETE_IMAGE",    // 64,活码删除图片
	"TIME_CLOCK_CREATE",            // 65,创建打卡
	"TIME_CLOCK_DELETE",            // 66,删除打卡
	"TIME_CLOCK_MEMBER_CREATE",     // 67,创建打卡成员
	"TIME_CLOCK_MEMBER_AUDIT",      // 68,审核打卡成员
	"TIME_CLOCK_MEMBER_DELETE",     // 69,删除打卡成员
	"TIME_CLOCK_UPDATE",            // 70,更新打卡
	"TIME_CLOCK_RECORD_CREATE",     // 71,创建打卡记录
	"TIME_CLOCK_RECORD_DELETE",     // 72,删除打卡记录
	"CDKEY_USE",                    // 73,使用Cdkey
	"",                             // 74,更新附件信息
	"",                             // 75,更新文章分类
	"",                             // 76,更新菜单
	"",                             // 77,更新文章标签
	"",                             // 78,更新站点信息
	"",                             // 79,更新主题
	"",                             // 80,更新收集表单
	"",                             // 81,更新友链
	"",                             // 82,更新系统变量
	"",                             // 83,更新图库图片
	"",                             // 84
	"",                             // 85
	"",                             // 86,修改文章
	"",                             // 87,修改自建页面
	"",                             // 88,修改主题配置
	"",                             // 89，修改主题文件
}

type UserLogType uint

// JsonDate反序列化
func (t *UserLogType) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range USER_LOG_TYPE {
		if status == value {
			*t = UserLogType(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t UserLogType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + USER_LOG_TYPE[t] + "\""), nil
}

func UserLogTypeValue(value string) (UserLogType, error) {
	for i, postType := range USER_LOG_TYPE {
		if postType == value {
			return UserLogType(i), nil
		}
	}
	return UserLogType(0), errors.New("未找到状态码：" + value)
}
