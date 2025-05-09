package repository

import (
	"github.com/CemAkan/url-shortener/config"
	"github.com/CemAkan/url-shortener/internal/domain"
	"gorm.io/gorm"
)

type URLRepository interface {
	Create(*domain.URL) error
	FindByCode(code string) (*domain.URL, error)
	FindByUserID(id uint) ([]domain.URL, error)
	Update(url *domain.URL) error
	IncrementTotalClicks(code string) error
}

type urlRepo struct {
	db *gorm.DB
}

func NewURLRepository() URLRepository {
	return &urlRepo{
		db: config.DB
	}
}

// Create inserts new url
func (r *urlRepo) Create (url *domain.URL)error{
	return r.db.Create(url).Error
}

// FindByCode retrieves URL by short code
func (r *urlRepo) FindByCode(code string) (*domain.URL, error){
	var url domain.URL
	err := r.db.Where("code = ?",code).First(&url).Error

	return &url, err
}
