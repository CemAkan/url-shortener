package app

import (
	"fmt"
	"github.com/CemAkan/url-shortener/config"
	"github.com/CemAkan/url-shortener/email"
	"github.com/CemAkan/url-shortener/internal/infrastructure/mail"
	"github.com/CemAkan/url-shortener/pkg/logger"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	mailVerificationMailSubject string = "Please verify your email by clicking the button below."
	passwordResetMailSubject    string = "Click the button below to reset your password. This link will expire in 15 minutes."
	mailLogger                  *logrus.Logger
)

func init() {
	mailLogger = logger.SpecialLogger("mail", "file")
}

type MailService interface {
	SendVerificationMail(name, baseUrl, emailAddr, verifyLink string) error
	SendPasswordResetMail(name, baseUrl, emailAddr, verifyLink string) error
	VerifyLinkGenerator(userID uint, baseURL, subject string, duration time.Duration) (string, error)
	GetMailLogger() *logrus.Logger
}

type mailService struct{}

// NewMailService constructs mailService struct
func NewMailService() MailService {
	return &mailService{}
}

// SendVerificationMail renders verification template and sends email
func (s *mailService) SendVerificationMail(name, baseUrl, emailAddr, verifyLink string) error {

	htmlBody, err := email.Render(email.EmailData{
		Title:            "Verify Your Email",
		Greeting:         fmt.Sprintf("Hello %s,", name),
		Message:          mailVerificationMailSubject,
		VerificationLink: verifyLink,
		LogoURL:          baseUrl + "/api/assets/logo.svg",
		HeaderURL:        baseUrl + "/api/assets/header.png",
		ButtonText:       "✔️ Verify Your Mail",
	})
	if err != nil {
		return err
	}

	return mail.Mail.Send(emailAddr, "Verify Your Email", htmlBody)
}

// SendPasswordResetMail renders reset-password template and sends email
func (s *mailService) SendPasswordResetMail(name, baseUrl, emailAddr, verifyLink string) error {
	htmlBody, err := email.Render(email.EmailData{
		Title:            "Reset Your Password",
		Greeting:         fmt.Sprintf("Hello %s,", name),
		Message:          passwordResetMailSubject,
		VerificationLink: verifyLink,
		LogoURL:          baseUrl + "/api/assets/logo.svg",
		HeaderURL:        baseUrl + "/api/assets/header.png",
		ButtonText:       "🔄 Reset Password",
	})
	if err != nil {
		return err
	}

	return mail.Mail.Send(emailAddr, "Reset Your Password", htmlBody)
}

// VerifyLinkGenerator generates tokenized link for verification or password reset
func (s *mailService) VerifyLinkGenerator(userID uint, baseURL, subject string, duration time.Duration) (string, error) {
	token, err := config.GenerateToken(userID, duration, subject)
	if err != nil {
		return "", err
	}

	return baseURL + "/" + token, nil
}

// GetMailLogger returns mail service logger
func (s *mailService) GetMailLogger() *logrus.Logger {
	return mailLogger
}
