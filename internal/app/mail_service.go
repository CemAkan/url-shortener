package app

import (
	"github.com/CemAkan/url-shortener/config"
	"github.com/CemAkan/url-shortener/internal/utils"
	"github.com/CemAkan/url-shortener/pkg/infrastructure"
	"time"
)

type MailService interface {
	SendVerificationMail(name, email, verifyLink string) error
	SendPasswordResetMail(name, email, verifyLink string) error
	VerifyLinkGenerator(userID uint, baseURL string) (string, error)
}

type mailService struct{}

// NewMailService constructs mailService struct
func NewMailService() MailService {
	return &mailService{}
}

// SendVerificationMail sends email to verify mail address
func (s *mailService) SendVerificationMail(name, email, verifyLink string) error {
	subject := "Welcome URL-Shortener! Please confirm your email address"
	if err := infrastructure.Mail.Send(email, subject, utils.GenerateEmailVerification(name, verifyLink)); err != nil {
		return err
	}
	return nil
}

// VerifyLinkGenerator generates new verification end of the link
func (s *mailService) VerifyLinkGenerator(userID uint, baseURL string) (string, error) {
	token, err := config.GenerateToken(userID, 24*time.Hour, "email_verification")
	if err != nil {
		return "", err
	}

	return baseURL + "/verify/" + token, nil

}

//

func (s *mailService) SendPasswordResetMail(name, email, verifyLink string) error {
	subject := " Reset your URL-Shortener password"

	if err := infrastructure.Mail.Send(email, subject, utils.GenerateResetPasswordEmail(name, verifyLink)); err != nil {
		return err
	}

	return nil
}
