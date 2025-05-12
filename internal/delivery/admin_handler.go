package delivery

import "github.com/CemAkan/url-shortener/internal/app"

type AdminHandler struct {
	userService app.UserService
	urlService  app.URLService
}

// NewAdminHandler constructor
func NewAdminHandler(userService app.UserService, urlService app.URLService) *AdminHandler {
	return &AdminHandler{
		userService: userService,
		urlService:  urlService,
	}
}
