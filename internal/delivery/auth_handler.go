package delivery

import "github.com/CemAkan/url-shortener/internal/app"

type AuthHandler struct {
	service app.UserService
}

// NewAuthHandler creates a new AuthHandler struct with given UserService input
func NewAuthHandler(userService app.UserService) *AuthHandler {
	return &AuthHandler{
		service: userService,
	}
}
