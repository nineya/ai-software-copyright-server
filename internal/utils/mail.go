package utils

import (
	"gopkg.in/gomail.v2"
	"mime"
	"path"
)

func SendMail() {
	m := gomail.NewMessage()
	m.SetHeader("From", mime.QEncoding.Encode("UTF-8", path.Base("玖涯博客 - 殇雪"))+"<blog@nineya.com>")
	m.SetHeader("To", "361654768@qq.com")
	m.SetHeader("Subject", "Hello 玖涯博客 - 殇雪!")
	m.SetBody("text/html", "Hello <b>Bob</b>玖涯博客 - 殇雪 and <i>Cora</i>!")

	d := gomail.NewDialer("smtp.ym.163.com", 994, "blog@nineya.com", "LSWang.361654768")
	d.SSL = true

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
