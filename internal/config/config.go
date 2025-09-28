package config

type Server struct {
	Port  int    `mapstructure:"port" json:"port" yaml:"port"`
	Mode  string `mapstructure:"mode" json:"mode" yaml:"mode"` // 运行模式，dev或prod
	Cache string `mapstructure:"cache" json:"cache" yaml:"cache"`
}

type Datasource struct {
	Type         string `mapstructure:"type" json:"type" yaml:"type"`
	Url          string `mapstructure:"url" json:"url" yaml:"url"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"maxIdleConns" yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"maxOpenConns" yaml:"max-open-conns"` // 打开到数据库的最大连接数
	LogLevel     string `mapstructure:"log-level" json:"logLevel" yaml:"log-level"`               // 日志
	UseCache     bool   `mapstructure:"use-cache" json:"useCache" yaml:"use-cache"`               // 是否开启缓存
	CacheSize    int    `mapstructure:"cache-size" json:"cacheSize" yaml:"cache-size"`            // 缓存数量
	ShowSql      bool   `mapstructure:"show-sql" json:"showSql" yaml:"show-sql"`                  // 显示sql
}

type Redis struct {
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`                   // redis的哪个数据库
	Host     string `mapstructure:"host" json:"host" yaml:"host"`             // 服务器地址:端口
	Password string `mapstructure:"password" json:"password" yaml:"password"` // 密码
}

type JWT struct {
	SigningKey  string `mapstructure:"signing-key" json:"signingKey" yaml:"signing-key"`    // jwt签名
	ExpiresTime int64  `mapstructure:"expires-time" json:"expiresTime" yaml:"expires-time"` // 过期时间
	BufferTime  int64  `mapstructure:"buffer-time" json:"bufferTime" yaml:"buffer-time"`    // 缓冲时间
}

type Weixin struct {
	Merchant    Merchant    `mapstructure:"merchant" json:"merchant" yaml:"merchant"`
	MiniProgram MiniProgram `mapstructure:"mini-program" json:"miniProgram" yaml:"mini-program"`
}

type Merchant struct {
	Mchid          string `mapstructure:"mchid" json:"mchid" yaml:"mchid"`
	ApiV3Key       string `mapstructure:"api-v3-key" json:"apiV3Key" yaml:"api-v3-key"`
	SerialNo       string `mapstructure:"serial-no" json:"serialNo" yaml:"serial-no"`
	PrivateKeyPath string `mapstructure:"private-key-path" json:"privateKeyPath" yaml:"private-key-path"`
	NotifyUrl      string `mapstructure:"notify-url" json:"notifyUrl" yaml:"notify-url"`
}

type MiniProgram struct {
	BloggerHelper   MiniProgramItem `mapstructure:"blogger-helper" json:"bloggerHelper" yaml:"blogger-helper"`
	TrafficToolbox  MiniProgramItem `mapstructure:"traffic-toolbox" json:"trafficToolbox" yaml:"traffic-toolbox"`
	OperationHelper MiniProgramItem `mapstructure:"operation-helper" json:"operationHelper" yaml:"operation-helper"`
	WechatToolbox   MiniProgramItem `mapstructure:"wechat-toolbox" json:"wechatToolbox" yaml:"wechat-toolbox"`
	NetdiskHelper   MiniProgramItem `mapstructure:"netdisk-helper" json:"netdiskHelper" yaml:"netdisk-helper"`
}

type MiniProgramItem struct {
	Appid  string `mapstructure:"appid" json:"appid" yaml:"appid"`
	Secret string `mapstructure:"secret" json:"secret" yaml:"secret"`
}

type ToolImage struct {
	AllowUnknownBucket bool     `mapstructure:"allow-unknown-bucket" json:"allowUnknownBucket" yaml:"allow-unknown-bucket"` // 允许上传时不指定bucket，使用默认bucket名称
	DefaultBucket      string   `mapstructure:"default-bucket" json:"defaultBucket" yaml:"default-bucket"`                  // 默认bucket名称
	Buckets            []string `mapstructure:"buckets" json:"buckets" yaml:"buckets"`                                      // 允许上传的bucket名称
	DownloadPath       string   `mapstructure:"download-path" json:"downloadPath" yaml:"download-path"`                     // 下载文件的域名路径前缀
}

type Tool struct {
	ToolImage ToolImage `mapstructure:"image" json:"image" yaml:"image"`
}

type PluginZhipuAi struct {
	Apikey string `mapstructure:"apikey" json:"apikey" yaml:"apikey"`
}

type PluginQuark struct {
	UserId      int64  `mapstructure:"user-id" json:"userId" yaml:"user-id"`
	Cookie      string `mapstructure:"cookie" json:"cookie" yaml:"cookie"`
	ToPdirFid   string `mapstructure:"to-pdir-fid" json:"toPdirFid" yaml:"to-pdir-fid"`
	AppendShare string `mapstructure:"append-share" json:"appendShare" yaml:"append-share"`
}

type PluginBaidu struct {
	UserId      int64  `mapstructure:"user-id" json:"userId" yaml:"user-id"`
	Cookie      string `mapstructure:"cookie" json:"cookie" yaml:"cookie"`
	ToPath      string `mapstructure:"to-path" json:"toPath" yaml:"to-path"`
	Pwd         string `mapstructure:"pwd" json:"pwd" yaml:"pwd"`
	AppendShare string `mapstructure:"append-share" json:"appendShare" yaml:"append-share"`
}

type PluginMail struct {
	Host      string `mapstructure:"host" json:"host" yaml:"host"`
	Port      int    `mapstructure:"port" json:"port" yaml:"port"`
	From      string `mapstructure:"from" json:"from" yaml:"from"`
	Username  string `mapstructure:"username" json:"username" yaml:"username"` //发件人邮箱账号
	Password  string `mapstructure:"password" json:"password" yaml:"password"`
	AdminMail string `mapstructure:"admin-mail" json:"adminMail" yaml:"admin-mail"`
}

type Plugin struct {
	ZhipuAi PluginZhipuAi `mapstructure:"zhipu-ai" json:"zhipuAi" yaml:"zhipu-ai"`
	Quark   PluginQuark   `mapstructure:"quark" json:"quark" yaml:"quark"`
	Baidu   PluginBaidu   `mapstructure:"baidu" json:"baidu" yaml:"baidu"`
	Mail    PluginMail    `mapstructure:"mail" json:"mail" yaml:"mail"`
}

type Logger struct {
	LogInConsole bool   `mapstructure:"log-in-console" json:"logInConsole" yaml:"log-in-console"` // 输出控制台
	Level        string `mapstructure:"level" json:"level" yaml:"level"`
}

type Config struct {
	Server     Server     `mapstructure:"server" json:"server" yaml:"server"`
	Datasource Datasource `mapstructure:"datasource" json:"datasource" yaml:"datasource"`
	Redis      Redis      `mapstructure:"redis" json:"redis" yaml:"redis"`
	JWT        JWT        `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Weixin     Weixin     `mapstructure:"weixin" json:"weixin" yaml:"weixin"`
	Tool       Tool       `mapstructure:"tool" json:"tool" yaml:"tool"`
	Plugin     Plugin     `mapstructure:"plugin" json:"plugin" yaml:"plugin"`
	Logger     Logger     `mapstructure:"logger" json:"logger" yaml:"logger"`
}
