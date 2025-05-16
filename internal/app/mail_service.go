package app

import (
	"fmt"
	"github.com/CemAkan/url-shortener/config"
	"github.com/CemAkan/url-shortener/email"
	"github.com/CemAkan/url-shortener/pkg/infrastructure"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	mailVerificationMailSubject string = "Please verify your email by clicking the button below."
	passwordResetMailSubject    string = "Click the button below to reset your password. This link will expire in 15 minutes."
	mailLogger                  *logrus.Logger
)

func init() {
	mailLogger = infrastructure.SpecialLogger("mail", "file")
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

	htmlBody, err := email.Render("verify-email", email.EmailData{
		Title:            "Verify Your Email",
		Greeting:         fmt.Sprintf("Hello %s,", name),
		Message:          mailVerificationMailSubject,
		VerificationLink: verifyLink,
		LogoURL:          baseUrl + "/assets/logo.svg",
		HeaderURL:        baseUrl + "/assets/header.png",
	})
	if err != nil {
		return err
	}

	return infrastructure.Mail.Send(emailAddr, mailVerificationMailSubject, htmlBody)
}

// SendPasswordResetMail renders reset-password template and sends email
func (s *mailService) SendPasswordResetMail(name, baseUrl, emailAddr, verifyLink string) error {
	htmlBody, err := email.Render("reset-password", email.EmailData{
		Title:            "Reset Your Password",
		Greeting:         fmt.Sprintf("Hello %s,", name),
		Message:          passwordResetMailSubject,
		VerificationLink: verifyLink,
		LogoURL:          baseUrl + "/assets/logo.svg",
		HeaderURL:        baseUrl + "/assets/header.png",
	})
	if err != nil {
		return err
	}

	return infrastructure.Mail.Send(emailAddr, passwordResetMailSubject, htmlBody)
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
