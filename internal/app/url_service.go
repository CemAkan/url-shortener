package app

import "github.com/CemAkan/url-shortener/internal/domain"

type URLService interface {
	Shorten(originalURL string, userID uint) (*domain.URL, error)
	GetByCode(code string) (*domain.URL, error)
	GetByUserID(userID uint) ([]domain.URL, error)
}
