package infrastructure

import (
	"crypto/tls"
	"github.com/CemAkan/url-shortener/config"
	"github.com/go-gomail/gomail"
	"strconv"
)

type Mailer struct {
	dialer *gomail.Dialer
	from   string
}

var Mail *Mailer

func InitMail() {

	port, _ := strconv.Atoi(config.GetEnv("SMTP_PORT", ""))
	host := config.GetEnv("SMTP_HOST", "")
	user := config.GetEnv("SMTP_USER", "")
	pass := config.GetEnv("SMTP_PASS", "")
	from := config.GetEnv("SMTP_FROM", "")

	dialer := gomail.NewDialer(host, port, user, pass)

	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	Mail = &Mailer{
		dialer: dialer,
		from:   from,
	}
}

// Send sends email
func (m *Mailer) Send(to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	return m.dialer.DialAndSend(msg)

}
