package response

type CaptchaResponse struct {
	CaptchaId string `json:"captchaId"`
	Path      string `json:"path"`
}
