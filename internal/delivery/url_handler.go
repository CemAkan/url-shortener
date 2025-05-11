package delivery

import "github.com/CemAkan/url-shortener/internal/app"

type URLhandler struct {
	service app.URLService
}

func NewURLService(urlService app.URLService) *URLhandler {
	return &URLhandler{
		service: urlService,
	}
}
