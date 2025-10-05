package enum

import (
	"github.com/pkg/errors"
)

var BUY_TYPE = [...]string{
	"",                          // 0,未定义
	"SOFTWARE_COPYRIGHT_CREATE", // 1,创建软著申请
	"REDBOOK_VALUATION",         // 2,小红书账号估值
	"REDBOOK_WEIGHT",            // 3,小红书账号权重检测
	"REDBOOK_REMOVE_WATERMARK",  // 4,小红书去水印
	"REDBOOK_TITLE",             // 5,小红书爆款标题生成
	"REDBOOK_NOTE",              // 6,小红书笔记帮写/优化
	"REDBOOK_PLANTING",          // 7,小红书种草笔记生成
	"SHORT_LINK_CLOUD_DISK",     // 8,短链网盘链接转换
	"QRCODE_BUILD",              // 9,二维码生成
	"SHORT_LINK_REDIRECT",       // 10,短链重定向
	"FLASH_PICTURE_ORIGIN",      // 11,查看闪照来源
	"FLASH_PICTURE_BROWSE",      // 12,浏览闪照
	"MP_IMAGETEXT_OPTIMIZE",     // 13,公众号图文帮写/优化
	"NETDISK_RESPONSE_CREATE",   // 14,创建网盘资源
	"REWARD_GOODS",              // 15,激励物品
	"NETDISK_RESPONSE_CHECK",    // 16,网盘资源检查
	"SHORT_LINK_STATISTIC",      // 17,短链分析服务
	"QRCODE_ADD_IMAGE",          // 18,活码添加图片
	"TIME_CLOCK_CREATE",         // 19,发起打卡活动
}

type BuyType uint

// JsonDate反序列化
func (t *BuyType) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range BUY_TYPE {
		if status == value {
			*t = BuyType(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t BuyType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + BUY_TYPE[t] + "\""), nil
}

func BuyTypeValue(value string) (BuyType, error) {
	for i, postType := range BUY_TYPE {
		if postType == value {
			return BuyType(i), nil
		}
	}
	return BuyType(0), errors.New("未找到状态码：" + value)
}
