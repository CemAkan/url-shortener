package response

import "github.com/CemAkan/url-shortener/internal/domain/model"

type UserURLsResponse struct {
	User model.User  `json:"user"`
	Urls []model.URL `json:"urls"`
}
