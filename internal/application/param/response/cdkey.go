package response

type CdkeyUseResponse struct {
	NyCredits int    `json:"nyCredits"`
	Remark    string `json:"remark,omitempty"`
}
