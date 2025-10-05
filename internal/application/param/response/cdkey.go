package response

type CdkeyUseResponse struct {
	Credits int    `json:"credits"`
	Remark  string `json:"remark,omitempty"`
}
