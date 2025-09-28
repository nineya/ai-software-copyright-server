package request

type WeixinCodeUnlimitParam struct {
	Page       string `json:"page"`
	Scene      string `json:"scene"` // 最大32个可见字符，只支持数字，大小写英文以及部分特殊字符：!#$&'()*+,/:;=?@-._~
	CheckPath  bool   `json:"check_path"`
	EnvVersion string `json:"env_version"` // 要打开的小程序版本。正式版为 "release"，体验版为 "trial"，开发版为 "develop"。默认是正式版。
	AutoColor  bool   `json:"auto_color"`  // 自动配置线条颜色，如果颜色依然是黑色，则说明不建议配置主色调，默认 false
	IsHyaline  bool   `json:"is_hyaline"`  // 默认是false，是否需要透明底色，为 true 时，生成透明底色的小程序
}

// https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/sec-center/sec-check/msgSecCheck.html
type WeixinMsgSecCheckParam struct {
	Content string `json:"content"` // 需检测的文本内容，文本字数的上限为2500字，需使用UTF-8编码
	Version int8   `json:"version"` // 接口版本号，2.0版本为固定值2
	Scene   int8   `json:"scene"`   // 场景枚举值（1 资料；2 评论；3 论坛；4 社交日志）
	Openid  string `json:"openid"`  // 用户的openid（用户需在近两小时访问过小程序）
}
