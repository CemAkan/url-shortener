package app

import (
	"github.com/CemAkan/url-shortener/config"
	"github.com/CemAkan/url-shortener/email"
	"github.com/CemAkan/url-shortener/pkg/infrastructure"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	mailVerificationMailSubject string = "Welcome URL-Shortener! Please confirm your email address"
	passwordResetMailSubject    string = "Reset your URL-Shortener password"
	mailLogger                  *logrus.Logger
)

func init() {
	mailLogger = infrastructure.SpecialLogger("mail", "file")
}

type MailService interface {
	SendVerificationMail(name, email, verifyLink string) error
	SendPasswordResetMail(name, email, verifyLink string) error
	VerifyLinkGenerator(userID uint, baseURL, subject string, duration time.Duration) (string, error)
	GetMailLogger() *logrus.Logger
}

type mailService struct{}

// NewMailService constructs mailService struct
func NewMailService() MailService {
	return &mailService{}
}

// SendVerificationMail renders verification template and sends email
func (s *mailService) SendVerificationMail(name, emailAddr, verifyLink string) error {
	htmlBody, err := email.Render("verify-email", email.EmailData{
		Title:            "Verify Your Email",
		VerificationLink: verifyLink,
	})
	if err != nil {
		return err
	}

	if err := infrastructure.Mail.Send(emailAddr, mailVerificationMailSubject, htmlBody); err != nil {
		return err
	}

	return nil
}

// SendPasswordResetMail renders reset-password template and sends email
func (s *mailService) SendPasswordResetMail(name, emailAddr, verifyLink string) error {
	htmlBody, err := email.Render("reset-password", email.EmailData{
		Title:            "Reset Your Password",
		VerificationLink: verifyLink,
	})
	if err != nil {
		return err
	}

	if err := infrastructure.Mail.Send(emailAddr, passwordResetMailSubject, htmlBody); err != nil {
		return err
	}

	return nil
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
