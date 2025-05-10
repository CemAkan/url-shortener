package app

import (
	"github.com/CemAkan/url-shortener/internal/domain"
	"github.com/CemAkan/url-shortener/internal/repository"
)

type URLService interface {
	Shorten(originalURL string, userID uint, customCode *string) (*domain.URL, error)
	GetByCode(code string) (*domain.URL, error)
	GetByUserID(userID uint) ([]domain.URL, error)
}

type urlService struct {
	repo repository.URLRepository
}

func NewURLService(urlRepo repository.URLRepository) URLService {
	return &userService{
		repo: urlRepo,
	}
}
