package common

type WelfareLabel struct {
	Label string        `json:"label"` // 标签名称
	Items []WelfareItem `json:"items"` // 羊毛选项
}

type WelfareItem struct {
	Title     string `json:"title"`     // 羊毛标题
	Desc      string `json:"desc"`      // 羊毛简介
	ImageUrl  string `json:"imageUrl"`  // 羊毛图片url
	TargetUrl string `json:"targetUrl"` // 羊毛目标地址
	InnerUrl  string `json:"innerUrl"`  // 羊毛内部地址
	WxAppid   string `json:"wxAppid"`   // 羊毛跳转小程序Appid
	WxPath    string `json:"wxPath"`    // 羊毛跳转小程序目标地址
}
