package repository

import "github.com/CemAkan/url-shortener/internal/domain"

type URLRepository interface {
	Create(*domain.URL) error
	FindByCode(code string) (*domain.URL, error)
	FindByUserID(id uint) ([]domain.URL, error)
	Update(url *domain.URL) error
	IncrementTotalClicks(code string) error
}
