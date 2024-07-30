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

type Logger struct {
	LogInConsole bool `mapstructure:"log-in-console" json:"logInConsole" yaml:"log-in-console"` // 输出控制台
}

type Config struct {
	Server     Server     `mapstructure:"server" json:"server" yaml:"server"`
	Datasource Datasource `mapstructure:"datasource" json:"datasource" yaml:"datasource"`
	Redis      Redis      `mapstructure:"redis" json:"redis" yaml:"redis"`
	JWT        JWT        `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Logger     Logger     `mapstructure:"logger" json:"logger" yaml:"logger"`
}
