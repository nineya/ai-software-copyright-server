package response

import "ai-software-copyright-server/internal/application/model/enum"

type RedbookProhibitedDetectionResponse struct {
	UserBuyResponse
	Content         string `json:"content"`
	ProhibitedCount int    `json:"prohibitedCount"` // 违禁词计数
	SensitiveCount  int    `json:"sensitiveCount"`  // 敏感词计数
	CustomCount     int    `json:"customCount"`     // 自定义词计数
}

type RedbookRemoveWatermarkResponse struct {
	UserBuyResponse
	Urls []string `json:"urls"`
}

type RedbookWriteTitleResponse struct {
	UserBuyResponse
	Titles []string `json:"titles"`
}

type RedbookWriteMessageResponse struct {
	UserBuyResponse
	Content string `json:"content"`
}

type RedbookValuationResponse struct {
	UserBuyResponse
	RedbookProfileInfoUserBasicInfo
	FollowCount      string `json:"followCount"`      // 关注数
	FansCount        string `json:"fansCount"`        // 粉丝数
	InteractionCount string `json:"interactionCount"` // 赞与收藏数
	Price            string `json:"price"`
}

type RedbookWeightResponse struct {
	UserBuyResponse
	Score       int                       `json:"score"`
	ScoreMsg    string                    `json:"scoreMsg"`
	Info        RedbookWeightInfoResponse `json:"info"`
	Nickname    RedbookWeightItemResponse `json:"nickname"`
	Fans        RedbookWeightItemResponse `json:"fans"`
	Interaction RedbookWeightItemResponse `json:"interaction"`
	Follow      RedbookWeightItemResponse `json:"follow"`
	Hot         RedbookWeightItemResponse `json:"hot"`
	Restrict    RedbookWeightItemResponse `json:"restrict"`
	Desc        RedbookWeightItemResponse `json:"desc"`
	Flow        RedbookWeightItemResponse `json:"flow"`
	PostTime    RedbookWeightItemResponse `json:"postTime"`
	LiveTime    RedbookWeightItemResponse `json:"liveTime"`
	Quality     RedbookWeightItemResponse `json:"quality"`
}

type RedbookWeightInfoResponse struct {
	RedbookProfileInfoUserBasicInfo
	FollowCount      string `json:"followCount"`      // 关注数
	FansCount        string `json:"fansCount"`        // 粉丝数
	InteractionCount string `json:"interactionCount"` // 赞与收藏数
}

type RedbookWeightItemResponse struct {
	Value string         `json:"value"`
	Hint  string         `json:"hint"`
	Level enum.HintLevel `json:"level"`
}

type RedbookProfileInfoResponse struct {
	User RedbookProfileInfoUser `json:"user"`
}

type RedbookProfileInfoUser struct {
	LoggedIn     bool                           `json:"loggedIn"` //是否登录
	UserPageData RedbookProfileInfoUserPageData `json:"userPageData"`
	Notes        [][]RedbookProfileInfoNoteItem `json:"notes"`
}

type RedbookProfileInfoNoteItem struct {
	Index    int                        `json:"index"`
	NoteCard RedbookProfileInfoNoteCard `json:"noteCard"`
}

type RedbookProfileInfoNoteCard struct {
	Type         string                                 `json:"type"`
	DisplayTitle string                                 `json:"displayTitle"`
	InteractInfo RedbookProfileInfoNoteCardInteractInfo `json:"interactInfo"`
}

type RedbookProfileInfoNoteCardInteractInfo struct {
	Sticky     bool   `json:"sticky"`
	Liked      bool   `json:"liked"`
	LikedCount string `json:"likedCount"`
}

type RedbookProfileInfoUserPageData struct {
	Interactions      []RedbookProfileInfoUserInteractionItem `json:"interactions"`
	Tags              []RedbookProfileInfoUserTagItem         `json:"tags"`
	UserAccountStatus RedbookProfileInfoAccountStatusInfo     `json:"userAccountStatus"`
	BasicInfo         RedbookProfileInfoUserBasicInfo         `json:"basicInfo"`
}

type RedbookProfileInfoUserInteractionItem struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Count string `json:"count"`
}

type RedbookProfileInfoUserTagItem struct {
	Name    string `json:"name"`
	TagType string `json:"tagType"`
}

type RedbookProfileInfoAccountStatusInfo struct {
	Type  int    `json:"type"`
	Toast string `json:"toast"`
}

type RedbookProfileInfoUserBasicInfo struct {
	Imageb     string `json:"imageb"`
	Nickname   string `json:"nickname"`
	RedId      string `json:"redId"`
	IpLocation string `json:"ipLocation"`
	Desc       string `json:"desc"`
}
