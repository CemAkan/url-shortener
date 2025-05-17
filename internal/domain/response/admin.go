package response

import "github.com/CemAkan/url-shortener/internal/domain/entity"

type UserURLsResponse struct {
	User entity.User  `json:"user"`
	Urls []entity.URL `json:"urls"`
}
