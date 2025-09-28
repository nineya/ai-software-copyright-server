package common

type SearchSiteConfig struct {
	Title        string                 `json:"title"`
	Subtitle     string                 `json:"subtitle"`
	Logo         string                 `json:"logo"`
	Favicon      string                 `json:"favicon"`
	Notice       string                 `json:"notice"`
	WechatQrcode string                 `json:"wechatQrcode"`
	Menus        []SearchSiteConfigLink `json:"menus"`
	Friends      []SearchSiteConfigLink `json:"friends"`
	Beian        string                 `json:"beian"`
	InlineCss    string                 `json:"inlineCss"`
	InlineJsBody string                 `json:"inlineJsBody"`
}

type SearchSiteConfigLink struct {
	Name string `json:"name"` // 菜单名称
	Url  string `json:"url"`
}
