package request

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"time"
)

type NetdiskCollectParam struct {
	Keyword string             `json:"keyword" form:"keyword" example:"壁纸" label:"搜索关键词"`
	Types   []enum.NetdiskType `json:"types" form:"types" label:"网盘类型列表"`
}

type NetdiskResourceQueryPageParam struct {
	QueryPageParam
	Type     string `json:"type" form:"type" label:"网盘类型"`
	Status   string `json:"status" form:"status" label:"资源状态"`
	Origin   string `json:"origin" form:"origin" label:"资源类型"`
	UserName string `json:"userName" form:"userName" label:"资源所属用户名称"`
}

type NetdiskResourceSearchParam struct {
	QueryPageParam
	SecureMode   string   `json:"secureMode" form:"secureMode" label:"所处的安全模式"`
	CollectTypes []string `json:"collectTypes" form:"collectTypes" label:"采集网盘类型列表"`
}

type NetdiskResourceTransformParam struct {
	Content string `json:"content" form:"content" label:"要转存的文案内容"` // 要转存转化的文案内容
}

type NetdiskResourceSaveParam struct {
	table.NetdiskResource
	QuarkTransferName string `json:"quarkTransferName" form:"quarkTransferName" example:"搜索转存" label:"夸克转存模板名称"`
	BaiduTransferName string `json:"baiduTransferName" form:"baiduTransferName" example:"搜索转存" label:"百度转存模板名称"`
}

type NetdiskResourceTextSaveParam struct {
	QuarkTransferName string `json:"quarkTransferName" form:"quarkTransferName" example:"搜索转存" label:"夸克转存模板名称"`
	BaiduTransferName string `json:"baiduTransferName" form:"baiduTransferName" example:"搜索转存" label:"百度转存模板名称"`
	CreateShortLink   bool   `json:"createShortLink" form:"createShortLink" label:"是否创建短链"`                               // 是否创建短链
	Content           string `json:"content" form:"content" example:"资源：https://pan.quark.cn/s/xxxxxxx" label:"要转存的文案内容"` // 要转存的文案内容
}

type NetdiskResourceAccountTransferParam struct {
	SourceUserName string `json:"sourceUserName" form:"sourceUserName" example:"源账号用户名" label:"源账号用户名"`
	SaveName       string `json:"saveName" form:"saveName" example:"搜索专用账号" label:"转存目标账号"`
	Limit          int    `json:"limit" form:"limit" label:"迁移数量限制"`
}

type NetdiskResourceFileListParam struct {
	AccountName string `json:"accountName" form:"accountName" example:"搜索专用账号" label:"夸克账号"`
	DirFid      string `json:"dirFid" form:"dirFid"`
	Page        int    `json:"page" form:"page" label:"当前页"`
}

type NetdiskResourceChangeAccountParam struct {
	UserId   int64  `json:"userId" form:"userId" label:"用户ID"`
	UserName string `json:"userName" form:"userName" binding:"required,lte=25" label:"源用户名"` // 源用户名
	//Cookie    string `json:"cookie" form:"cookie" binding:"required,lte=1000"`      // 目标cookie
	ToPdirFid string `json:"toPdirFid" form:"toPdirFid" binding:"required,lte=100" label:"目标目录pid"` // 目标目录pid
	Limit     int    `json:"limit" form:"limit" label:"迁移数量限制"`
}

type NetdiskHelperUpdateExpireTime struct {
	ExpireTime       time.Time `json:"expireTime,omitempty" form:"expireTime" label:"过期时间"`
	WechatExpireTime time.Time `json:"wechatExpireTime,omitempty" form:"expireTime" label:"微信工具人过期时间"`
}

type NetdiskResourceUpdateInBatchParam struct {
	BatchParam
	Status string `json:"status" form:"status" label:"资源状态"`
}

type NetdiskQuarkShareParam struct {
	ExpiredType int      `json:"expired_type,omitempty" form:"expired_type,omitempty"`
	UrlType     int      `json:"url_type,omitempty" form:"url_type,omitempty"`
	FidList     []string `json:"fid_list,omitempty" form:"fid_list,omitempty"`
	Title       string   `json:"title,omitempty" form:"title,omitempty"`
}

type NetdiskQuarkGetShareLinkParam struct {
	ShareId string `json:"share_id,omitempty" form:"share_id,omitempty"`
}

type NetdiskQuarkGetStokenParam struct {
	PwdId    string `json:"pwd_id,omitempty" form:"pwd_id,omitempty"`
	Passcode string `json:"passcode,omitempty" form:"passcode,omitempty"`
}

type NetdiskQuarkSaveParam struct {
	FidList      []string `json:"fid_list,omitempty" form:"fid_list,omitempty"`             // 需要保存的文件id
	FidTokenList []string `json:"fid_token_list,omitempty" form:"fid_token_list,omitempty"` // 需要保存的文件token
	PdirFid      string   `json:"pdir_fid,omitempty" form:"pdir_fid,omitempty"`
	PwdId        string   `json:"pwd_id,omitempty" form:"pwd_id,omitempty"` // 夸克网盘分享链接的pwd_id
	Scene        string   `json:"scene,omitempty" form:"scene,omitempty"`
	Stoken       string   `json:"stoken,omitempty" form:"stoken,omitempty"`           // 夸克网盘分享链接的token
	ToPdirFid    string   `json:"to_pdir_fid,omitempty" form:"to_pdir_fid,omitempty"` //保存目标目录
}

type NetdiskQuarkNewDirParam struct {
	DirInitLock bool   `json:"dir_init_lock,omitempty" form:"dir_init_lock,omitempty"`
	DirPath     string `json:"dir_path,omitempty" form:"dir_path,omitempty"`
	FileName    string `json:"file_name,omitempty" form:"file_name,omitempty"` //新建文件夹名称
	PdirFid     string `json:"pdir_fid,omitempty" form:"pdir_fid,omitempty"`   // 保存文件夹的上级文件夹ID
}

type NetdiskQuarkDeleteFileParam struct {
	ActionType  int      `json:"action_type,omitempty" form:"action_type,omitempty"`
	ExcludeFids []string `json:"exclude_fids,omitempty" form:"exclude_fids,omitempty"`
	Filelist    []string `json:"filelist,omitempty" form:"filelist,omitempty"`
}

type NetdiskQuarkDeleteShareParam struct {
	ShareIds []string `json:"share_ids,omitempty" form:"share_ids,omitempty"`
}
