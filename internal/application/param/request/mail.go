package request

type MailParam struct {
	To      string `json:"to" form:"to"`
	Subject string `json:"subject" form:"subject"`
	Content string `json:"content" form:"content"`
}
