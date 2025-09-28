package common

import "ai-software-copyright-server/internal/application/model/enum"

type NetdiskHelperConfigQuark struct {
	Accounts         []NetdiskHelperConfigQuarkAccount  `json:"accounts,omitempty"`  // 夸克账号列表
	Transfers        []NetdiskHelperConfigQuarkTransfer `json:"transfers,omitempty"` // 转存模板列表
	SearchName       string                             `json:"searchName"`          // 搜索保存用的夸克名字
	DeleteSearchTime int                                `json:"deleteSearchTime"`    // 自动删除超过指定分钟的搜索转存
}

type NetdiskHelperConfigQuarkAccount struct {
	Name          string `json:"name"`          // 夸克名字
	ToPdirFid     string `json:"toPdirFid"`     // TODO, 转存目标目录fid
	AppendShare   string `json:"appendShare"`   // TODO, 追加转存
	AutoSyncShare bool   `json:"autoSyncShare"` // 自动同步分享
	Cookie        string `json:"cookie"`        // cookie
}

type NetdiskHelperConfigQuarkTransfer struct {
	Name        string `json:"name"`        // 转存模板名字
	AccountName string `json:"accountName"` // 转存目标账号名称
	ToPdirFid   string `json:"toPdirFid"`   // 转存目标目录fid
	AppendShare string `json:"appendShare"` // 追加转存
}

// 组合完整的夸克转存配置信息
type NetdiskHelperConfigQuarkTransferInfo struct {
	NetdiskHelperConfigQuarkTransfer
	Cookie string `json:"cookie"` // cookie
}

type NetdiskHelperConfigBaidu struct {
	Accounts         []NetdiskHelperConfigBaiduAccount  `json:"accounts,omitempty"`  // 百度账号列表
	Transfers        []NetdiskHelperConfigBaiduTransfer `json:"transfers,omitempty"` // 百度模板列表
	SearchName       string                             `json:"searchName"`          // 搜索保存用的夸克名字
	DeleteSearchTime int                                `json:"deleteSearchTime"`    // 自动删除超过指定分钟的搜索转存
}

type NetdiskHelperConfigBaiduAccount struct {
	Name          string `json:"name"`          // 夸克名字
	AutoSyncShare bool   `json:"autoSyncShare"` // 自动同步分享
	Cookie        string `json:"cookie"`        // cookie
}

type NetdiskHelperConfigBaiduTransfer struct {
	Name        string `json:"name"`        // 转存模板名字
	AccountName string `json:"accountName"` // 转存目标账号名称
	ToPath      string `json:"toPath"`      // 转存目标目录
	AppendShare string `json:"appendShare"` // 追加转存内容
	Pwd         string `json:"pwd"`         // 分享密码
}

// 组合完整的百度转存配置信息
type NetdiskHelperConfigBaiduTransferInfo struct {
	NetdiskHelperConfigBaiduTransfer
	Cookie string `json:"cookie"` // cookie
}

type NetdiskHelperConfigAi struct {
	ZhipuApiKey string `json:"zhipuApiKey"` // 智谱ApiKey
}

type NetdiskHelperConfigMail struct {
	Templates []NetdiskHelperConfigMailTemplate `json:"templates,omitempty"` // 邮件列表
	Host      string                            `json:"host"`
	Port      int                               `json:"port"`
	From      string                            `json:"from"`
	Username  string                            `json:"username"`
	Password  string                            `json:"password"`
}

type NetdiskHelperConfigMailTemplate struct {
	Name    string `json:"name,omitempty"` // 邮件名
	Subject string `json:"subject"`        // 邮件标题
	Content string `json:"content"`        // 邮件内容
}

type NetdiskHelperConfigWechat struct {
	HttpApi                    string                             `json:"httpApi"`   // 微信机器人的HttpApi
	AdminWxid                  string                             `json:"adminWxid"` // 管理微信wxid
	CollectTypes               []enum.NetdiskType                 `json:"collectTypes" label:"采集网盘类型列表"`
	SearchResultTemplate       string                             `json:"searchResultTemplate"`       // 搜索结果文案模板
	SearchFailTemplate         string                             `json:"searchFailTemplate"`         // 搜索失败文案模板
	SearchResourceTemplate     string                             `json:"searchResourceTemplate"`     // 搜索取得资源文案模板
	SearchResourceFailTemplate string                             `json:"searchResourceFailTemplate"` // 搜索取得资源失败文案模板
	Groups                     []NetdiskHelperConfigWechatGroup   `json:"groups,omitempty"`           // 微信群
	Tasks                      []NetdiskHelperConfigWechatTask    `json:"tasks,omitempty"`            // 定时任务
	Relays                     []NetdiskHelperConfigWechatRelay   `json:"relays,omitempty"`           // 转发任务
	Monitors                   []NetdiskHelperConfigWechatMonitor `json:"monitors,omitempty"`         // 监控群
}

type NetdiskHelperConfigWechatGroup struct {
	Name       string   `json:"name"`       // 群名字
	Tags       []string `json:"tags"`       // 群标签
	Wxid       string   `json:"wxid"`       // 群wxid
	WelcomeMsg string   `json:"welcomeMsg"` // 欢迎语
}

type NetdiskHelperConfigWechatTask struct {
	Name      string              `json:"name"`      // 任务名字
	Cron      string              `json:"cron"`      // 定时规则
	GroupTags []string            `json:"groupTags"` // 通知的群标签
	Type      enum.ClientTaskType `json:"type"`      // 任务类型
	Message   string              `json:"message"`   // 任务消息内容
}

type NetdiskHelperConfigWechatRelay struct { // 群转发
	Name            string   `json:"name"`            // 任务名字
	GroupName       string   `json:"groupName"`       // 需要转发的群名称
	MemberWxids     []string `json:"memberWxids"`     // 需要转发的成员wxid
	TriggerKeywords []string `json:"triggerKeywords"` // 触发关键字
	ExcludeKeywords []string `json:"excludeKeywords"` // 需要排除的关键词
	RelayTags       []string `json:"relayTags"`       // 转发目标群标签
}

type NetdiskHelperConfigWechatMonitor struct { // 群监控转发
	Name               string   `json:"name"`               // 监控名字
	GroupTags          []string `json:"groupTags"`          // 监控群标签
	ExcludeMemberWxids []string `json:"excludeMemberWxids"` // 白名单成员wxid
	TriggerKeywords    []string `json:"triggerKeywords"`    // 触发关键字
	CheckPrompt        string   `json:"checkPrompt"`        // AI检查提示词
	WarnMsg            string   `json:"warnMsg"`            // 警告消息
}
