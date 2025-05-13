package config

import (
	"github.com/CemAkan/url-shortener/pkg/infrastructure"
	"github.com/wneessen/go-mail"
	"strconv"
)

var MailClient *mail.Client

func InitMail() {

	port, _ := strconv.Atoi(GetEnv("SMTP_PORT", ""))

	client, err := mail.NewClient(
		GetEnv("SMTP_HOST", ""),
		mail.WithPort(port),
		mail.WithUsername(GetEnv("SMTP_USER", "")),
		mail.WithPassword(GetEnv("SMTP_PASS", "")),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
	)
	if err != nil {
		infrastructure.Log.Fatalf("Failed to initialize mail client: %v", err)
	}
	MailClient = client
}
