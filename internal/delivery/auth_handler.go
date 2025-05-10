package delivery

import "github.com/CemAkan/url-shortener/internal/app"

type AuthHandler struct {
	service app.UserService
}

func NewAuthHandler(userService app.UserService) *AuthHandler {
	return &AuthHandler{
		service: userService
	}
}
