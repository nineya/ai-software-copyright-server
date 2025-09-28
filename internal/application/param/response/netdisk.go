package response

type NetdiskHelperClientResponse struct {
	Online                bool   `json:"online" form:"online"`
	Version               string `json:"version" form:"version"`
	IsUpdate              bool   `json:"isUpdate" form:"isUpdate"`
	NewVersion            string `json:"newVersion" form:"newVersion"`
	NewVersionDownloadUrl string `json:"newVersionDownloadUrl" form:"newVersionDownloadUrl"`
	NewVersionDescription string `json:"newVersionDescription" form:"newVersionDescription"`
}

type NetdiskResourceCreateResponse struct {
	UserBuyResponse
	IsUpdate bool `json:"isUpdate" form:"isUpdate"`
}

type NetdiskResourceCheckResponse struct {
	UserBuyResponse
	Page *PageResponse `json:"page" form:"page"`
}

type NetdiskResourceTextSaveResponse struct {
	Content  string   `json:"content" form:"content"` // 要转存的文案内容
	FailUrls []string `json:"failUrls" form:"failUrls"`
}

type NetdiskQuarkResponse[T any] struct {
	Status    int    `json:"status" form:"status"`
	Code      int    `json:"code" form:"code"`
	Message   string `json:"message" form:"message"`
	Timestamp int64  `json:"timestamp" form:"timestamp"`
	Data      T      `json:"data" form:"data"`
	Metadata  any    `json:"metadata" form:"metadata"`
}

// 查询stoken
type NetdiskQuarkGetStokenData struct {
	Stoken    string                          `json:"stoken,omitempty" form:"stoken,omitempty"`
	ShareType int                             `json:"share_type" form:"share_type"`
	Title     string                          `json:"title,omitempty" form:"title,omitempty"`
	Author    NetdiskQuarkGetStokenDataAuthor `json:"author" form:"author"`
}

type NetdiskQuarkGetStokenDataAuthor struct {
	MemberType string `json:"member_type,omitempty" form:"member_type,omitempty"`
	NickName   string `json:"nick_name,omitempty" form:"nick_name,omitempty"`
	AvatarUrl  string `json:"avatar_url,omitempty" form:"avatar_url,omitempty"`
}

// 转存或分享响应的任务
type NetdiskQuarkTaskIdData struct {
	TaskId   string `json:"task_id" form:"task_id"`
	TaskSync bool   `json:"task_sync" form:"task_sync"`
}

// 查询分享链接
type NetdiskQuarkDetailData struct {
	IsOwner int                          `json:"is_owner" form:"is_owner"`
	List    []NetdiskQuarkDetailDataList `json:"list" form:"list"`
	Share   NetdiskQuarkDetailDataShare  `json:"share" form:"share"`
}

type NetdiskQuarkDetailDataList struct {
	Fid           string `json:"fid" form:"fid"`
	FileName      string `json:"file_name" form:"file_name"`
	ShareFidToken string `json:"share_fid_token" form:"share_fid_token"`
}

type NetdiskQuarkDetailDataShare struct {
	Title    string `json:"title" form:"title"`
	FileNum  int    `json:"file_num" form:"file_num"`
	FirstFid string `json:"first_fid" form:"first_fid"`
	ShareId  string `json:"share_id" form:"share_id"`
	ShareUrl string `json:"share_url" form:"share_url"`
}

// 查询文件列表链接
type NetdiskQuarkFileSortData struct {
	List []NetdiskQuarkFileSortDataList `json:"list" form:"list"`
}

type NetdiskQuarkFileSortDataList struct {
	Fid      string `json:"fid" form:"fid"`
	PdirFid  string `json:"pdir_fid" form:"pdir_fid"` // 当前所处目录的fid
	FileName string `json:"file_name" form:"file_name"`
	Dir      bool   `json:"dir" form:"dir"`
}

// 查询任务结果
type NetdiskQuarkTaskData struct {
	TaskId   string                     `json:"task_id" form:"task_id"`
	ShareId  string                     `json:"share_id" form:"share_id"`
	Status   int                        `json:"status" form:"status"`
	TaskType int                        `json:"task_type" form:"task_type"`
	SaveAs   NetdiskQuarkTaskDataSaveAs `json:"save_as" form:"save_as"`
}

type NetdiskQuarkTaskDataSaveAs struct {
	SaveAsSumNum  int      `json:"save_as_sum_num" form:"save_as_sum_num"`   // 转存的文件数量
	ToPdirFid     string   `json:"to_pdir_fid" form:"to_pdir_fid"`           // 转存的目标目录
	SaveAsTopFids []string `json:"save_as_top_fids" form:"save_as_top_fids"` // 文件/目录转存后的新fids
}

// -1 => '链接错误，链接失效或缺少提取码或访问频繁风控',
// -4 => '无效登录。请退出账号在其他地方的登录',
// -6 => '请用浏览器无痕模式获取 Cookie 后再试',
// -7 => '转存失败，转存文件夹名有非法字符，不能包含 < > | * ? \\ :，请改正目录名后重试'
// -8 => '转存失败，目录中已有同名文件或文件夹存在',
// -9 => '链接不存在或提取码错误',
// -10 => '转存失败，容量不足',
// -12 => '链接错误，提取码错误',
// -62 => '链接访问次数过多，请手动转存或稍后再试',
// 0 => '转存成功',
// 2 => '转存失败，目标目录不存在',
// 4 => '转存失败，目录中存在同名文件',
// 12 => '转存失败，转存文件数超过限制',
// 20 => '转存失败，容量不足',
// 105 => '链接错误，所访问的页面不存在',
// 115 => '该文件禁止分享'
// 117 => 分享文件已过期
// 145 => '分享已取消/被删除'
type NetdiskBaiduResponse struct {
	Errno  int    `json:"errno" form:"errno"`
	ErrMsg string `json:"err_msg" form:"err_msg"`
}

type NetdiskBaiduVerifyPassCodeResponse struct {
	NetdiskBaiduResponse
	Randsk string `json:"randsk" form:"randsk"`
}

type NetdiskBaiduTaskResponse struct {
	NetdiskBaiduResponse
	TaskId int64 `json:"taskid" form:"taskid"`
}

type NetdiskBaiduTransferResponse struct {
	NetdiskBaiduResponse
	Extra NetdiskBaiduTransferExtra  `json:"extra" form:"info"`
	Info  []NetdiskBaiduTransferInfo `json:"info" form:"info"`
}

type NetdiskBaiduTransferExtra struct {
	List []NetdiskBaiduTransferExtraItem `json:"list" form:"list"`
}

type NetdiskBaiduTransferExtraItem struct {
	From     string `json:"from" form:"from"` //源文件路径
	FromFsId int64  `json:"from_fs_id" form:"from_fs_id"`
	To       string `json:"to" form:"to"` //源文件路径
	ToFsId   int64  `json:"to_fs_id" form:"to_fs_id"`
}

type NetdiskBaiduTransferInfo struct {
	Errno int    `json:"errno" form:"errno"`
	FsId  int64  `json:"fsid" form:"fsid"`
	Path  string `json:"path" form:"path"`
}

type NetdiskBaiduNewDirResponse struct {
	NetdiskBaiduResponse
	FsId  int64  `json:"fs_id" form:"fs_id"`
	IsDir int    `json:"isdir" form:"isdir"` // 1表示目录
	Name  string `json:"name" form:"name"`
	Path  string `json:"path" form:"path"`
}

type NetdiskBaiduShareResponse struct {
	NetdiskBaiduResponse
	ShareUrl string //自己加的，保存带密码的分享链接
	ShortUrl string `json:"shorturl" form:"shorturl"`
	Link     string `json:"link" form:"link"`
}

type NetdiskBaiduListResponse struct {
	NetdiskBaiduResponse
	List []NetdiskBaiduListItem `json:"list" form:"list"`
}

type NetdiskBaiduListItem struct {
	Md5            string `json:"md5" form:"md5"` // 如果是文件有md5
	Path           string `json:"path" form:"path"`
	ServerFilename string `json:"server_filename" form:"server_filename"`
	FsId           int64  `json:"fs_id" form:"fs_id"`
	IsDir          int    `json:"isdir" form:"isdir"` // 1表示目录
}

type NetdiskBaiduShareInfoByHtmlResponse struct {
	NetdiskBaiduResponse
	ShareId      int64                                 `json:"shareid" form:"shareid"` // 分享id
	ShareUk      string                                `json:"share_uk" form:"share_uk"`
	LinkUserName string                                `json:"linkusername" form:"linkusername"`
	FileList     []NetdiskBaiduShareInfoByHtmlFileItem `json:"file_list" form:"file_list"`
}
type NetdiskBaiduShareInfoByHtmlFileItem struct {
	FsId           int64  `json:"fs_id" form:"fs_id"` // 资源文件id
	ServerFilename string `json:"server_filename" form:"server_filename"`
}

type NetdiskBaiduShareInfoResponse struct {
	NetdiskBaiduResponse
	ShareId int64                       `json:"share_id" form:"share_id"` // 分享id
	Uk      int                         `json:"uk" form:"uk"`             // 同share_uk
	List    []NetdiskBaiduShareInfoItem `json:"list" form:"list"`
}
type NetdiskBaiduShareInfoItem struct {
	FsId           string `json:"fs_id" form:"fs_id"` // 资源文件id
	ServerFilename string `json:"server_filename" form:"server_filename"`
	Path           string `json:"path" form:"path"`
}

type NetdiskBaiduGetBdstokenResponse struct {
	NetdiskBaiduResponse
	Result NetdiskBaiduGetBdstokenResult `json:"result" form:"result"`
}

type NetdiskBaiduGetBdstokenResult struct {
	Bdstoken  string `json:"bdstoken" form:"bdstoken"`
	Token     string `json:"token" form:"token"`
	Uk        int    `json:"uk" form:"uk"`
	IsDocUser int    `json:"isdocuser" form:"isdocuser"`
}
