package common

import "ai-software-copyright-server/internal/application/model/enum"

type ClientConfig struct {
	Netdisk ClientConfigNetdisk `json:"netdisk"`
	Mail    ClientConfigMail    `json:"mail"`
	Wechat  ClientConfigWechat  `json:"wechat"`
	QQ      ClientConfigQQ      `json:"qq"`
}

type ClientConfigNetdisk struct {
	Quarks           []ClientConfigNetdiskQuark `json:"quarks,omitempty"` // 夸克列表
	SearchQuarkName  string                     `json:"searchQuarkName"`  // 搜索保存用的夸克名字
	DeleteSearchTime int                        `json:"deleteSearchTime"` // 自动删除超过指定分钟的搜索转存
}

type ClientConfigNetdiskQuark struct {
	Name          string `json:"name"`          // 夸克名字
	ToPdirFid     string `json:"toPdirFid"`     // 转存目标目录fid
	AppendShare   string `json:"appendShare"`   // 追加转存
	AutoSyncShare bool   `json:"autoSyncShare"` // 自动同步分享
	Cookie        string `json:"cookie"`        // cookie
}

type ClientConfigMail struct {
	Templates []ClientConfigMailTemplate `json:"templates,omitempty"` // 邮件列表
	Host      string                     `json:"host"`
	Port      int                        `json:"port"`
	From      string                     `json:"from"`
	Username  string                     `json:"username"`
	Password  string                     `json:"password"`
}

type ClientConfigMailTemplate struct {
	Name    string `json:"name,omitempty"` // 邮件名
	Subject string `json:"subject"`        // 邮件标题
	Content string `json:"content"`        // 邮件内容
}

type ClientConfigWechat struct {
	HttpApi                    string                    `json:"httpApi"`                    // 微信机器人的HttpApi
	AdminWxid                  string                    `json:"adminWxid"`                  // 管理微信wxid
	Groups                     []ClientConfigWechatGroup `json:"groups,omitempty"`           // 微信群
	Tasks                      []ClientConfigWechatTask  `json:"tasks,omitempty"`            // 定时任务
	SearchGroupTags            []string                  `json:"searchGroupTags,omitempty"`  // 搜索用的群标签
	SearchResultTemplate       string                    `json:"searchResultTemplate"`       // 搜索结果文案模板
	SearchFailTemplate         string                    `json:"searchFailTemplate"`         // 搜索失败文案模板
	SearchResourceTemplate     string                    `json:"searchResourceTemplate"`     // 搜索取得资源文案模板
	SearchResourceFailTemplate string                    `json:"searchResourceFailTemplate"` // 搜索取得资源失败文案模板
}

type ClientConfigWechatGroup struct {
	Name       string   `json:"name"`       // 群名字
	Tags       []string `json:"tags"`       // 群标签
	Wxid       string   `json:"wxid"`       // 群wxid
	WelcomeMsg string   `json:"welcomeMsg"` // 欢迎语
}

type ClientConfigWechatTask struct {
	Name      string              `json:"name"`      // 任务名字
	Cron      string              `json:"cron"`      // 定时规则
	GroupTags []string            `json:"groupTags"` // 通知的群标签
	Type      enum.ClientTaskType `json:"type"`      // 任务类型
	Message   string              `json:"message"`   // 任务消息内容
}

type ClientConfigQQ struct {
	HttpApi                    string                `json:"httpApi"`                    // 微信机器人的HttpApi
	Groups                     []ClientConfigQQGroup `json:"groups,omitempty"`           // 微信群
	Tasks                      []ClientConfigQQTask  `json:"tasks,omitempty"`            // 定时任务
	SearchGroupTags            []string              `json:"searchGroupTags,omitempty"`  // 搜索用的群标签
	SearchResultTemplate       string                `json:"searchResultTemplate"`       // 搜索结果文案模板
	SearchFailTemplate         string                `json:"searchFailTemplate"`         // 搜索失败文案模板
	SearchResourceTemplate     string                `json:"searchResourceTemplate"`     // 搜索取得资源文案模板
	SearchResourceFailTemplate string                `json:"searchResourceFailTemplate"` // 搜索取得资源失败文案模板
}

type ClientConfigQQGroup struct {
	Name            string   `json:"name"`                      // 群名字
	Tags            []string `json:"tags"`                      // 群标签
	GroupId         string   `json:"groupId"`                   // 群号
	WelcomeMsg      string   `json:"welcomeMsg,omitempty"`      // 欢迎语
	WelcomeMailName string   `json:"welcomeMailName,omitempty"` // 欢迎邮件名
}

type ClientConfigQQTask struct {
	Name      string   `json:"name"`      // 任务名字
	Cron      string   `json:"cron"`      // 定时规则
	GroupTags []string `json:"groupTags"` // 通知的群标签
	Message   string   `json:"message"`   // 消息内容
}
