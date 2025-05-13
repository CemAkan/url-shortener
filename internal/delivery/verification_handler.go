package delivery

import "github.com/CemAkan/url-shortener/internal/app"

type VerificationHandler struct {
	userService app.UserService
	mailService app.MailService
}

// NewVerificationHandler generate a new VerificationHandler struct with given UserService and mailService inputs
func NewVerificationHandler(userService app.UserService, mailService app.MailService) *VerificationHandler {
	return &VerificationHandler{
		userService: userService,
		mailService: mailService,
	}
}
