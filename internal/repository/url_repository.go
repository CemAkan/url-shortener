package repository

import (
	"github.com/CemAkan/url-shortener/internal/domain"
	"github.com/CemAkan/url-shortener/pkg/infrastructure"
	"gorm.io/gorm"
)

type URLRepository interface {
	Create(*domain.URL) error
	FindByCode(code string) (*domain.URL, error)
	FindByUserID(id uint) ([]domain.URL, error)
	Update(url *domain.URL) error
	AddToTotalClicks(code string, count int) error
	Delete(code string) error
	DeleteUserAllUrls(userID uint) ([]domain.URL, error)
}

type urlRepo struct {
	db *gorm.DB
}

func NewURLRepository() URLRepository {
	return &urlRepo{
		db: infrastructure.DB,
	}
}

// Create inserts new url
func (r *urlRepo) Create(url *domain.URL) error {
	return r.db.Create(url).Error
}

// FindByCode retrieves URL by short code
func (r *urlRepo) FindByCode(code string) (*domain.URL, error) {
	var url domain.URL
	err := r.db.Where("code = ?", code).First(&url).Error

	return &url, err
}

// FindByUserID retrieves all URLs which associated with UserID
func (r *urlRepo) FindByUserID(id uint) ([]domain.URL, error) {
	var urls []domain.URL

	err := r.db.Where("user_id = ?", id).Find(&urls).Error

	return urls, err
}

// Update modifies to existing url
func (r *urlRepo) Update(url *domain.URL) error {
	return r.db.Save(url).Error
}

// AddToTotalClicks adds wanted click count to total clicks
func (r *urlRepo) AddToTotalClicks(code string, count int) error {
	return r.db.Model(&domain.URL{}).Where("code = ?", code).UpdateColumn("total_clicks", gorm.Expr("total_clicks + ?", count)).Error
}

// Delete removes code related url record from database
func (r *urlRepo) Delete(code string) error {
	if err := r.db.Where("code= ?", code).Delete(&domain.URL{}).Error; err != nil {
		return err
	}
	return nil
}

// DeleteUserAllUrls removes all userId related url records from database
func (r *urlRepo) DeleteUserAllUrls(userID uint) ([]domain.URL, error) {
	var urls []domain.URL

	if err := r.db.Where("user_id = ?", userID).Find(&urls).Error; err != nil {
		return nil, err
	}

	if err := r.db.Where("user_id = ?", userID).Delete(&domain.URL{}).Error; err != nil {
		return nil, err
	}

	return urls, nil
}
