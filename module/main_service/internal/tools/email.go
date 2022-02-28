package tools

import (
	"net/smtp"
	"time"

	"frog/module/main_service/internal/config"

	"github.com/jordan-wright/email"
)

var (
	auth smtp.Auth
	pool *email.Pool
)

func init() {
	auth = smtp.PlainAuth("", config.GetEmailConfig().Username, config.GetEmailConfig().Password, config.GetEmailConfig().EmailHost)
	pool, _ = email.NewPool(config.GetEmailConfig().EmailHost+":25", 1, auth)
}

func SendEmail(to string, content string) error {
	e := email.NewEmail()
	e.From = config.GetEmailConfig().Username
	e.To = []string{to}
	e.Subject = "Frog 验证码"
	e.Text = []byte(content)
	return pool.Send(e, 5*time.Second)
}
