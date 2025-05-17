package repository

import (
	"github.com/CemAkan/url-shortener/internal/domain/model"
	"github.com/CemAkan/url-shortener/pkg/infrastructure/db"
	"gorm.io/gorm"
)

type URLRepository interface {
	Create(*model.URL) error
	FindByCode(code string) (*model.URL, error)
	FindByUserID(id uint) ([]model.URL, error)
	Update(url *model.URL) error
	AddToTotalClicks(code string, count int) error
	Delete(code string) error
	DeleteUserAllUrls(userID uint) ([]model.URL, error)
}

type urlRepo struct {
	db *gorm.DB
}

func NewURLRepository() URLRepository {
	return &urlRepo{
		db: db.DB,
	}
}

// Create inserts new url
func (r *urlRepo) Create(url *model.URL) error {
	return r.db.Create(url).Error
}

// FindByCode retrieves URL by short code
func (r *urlRepo) FindByCode(code string) (*model.URL, error) {
	var url model.URL
	err := r.db.Where("code = ?", code).First(&url).Error

	return &url, err
}

// FindByUserID retrieves all URLs which associated with UserID
func (r *urlRepo) FindByUserID(id uint) ([]model.URL, error) {
	var urls []model.URL

	err := r.db.Where("user_id = ?", id).Find(&urls).Error

	return urls, err
}

// Update modifies to existing url
func (r *urlRepo) Update(url *model.URL) error {
	return r.db.Save(url).Error
}

// AddToTotalClicks adds wanted click count to total clicks
func (r *urlRepo) AddToTotalClicks(code string, count int) error {
	return r.db.Model(&model.URL{}).Where("code = ?", code).UpdateColumn("total_clicks", gorm.Expr("total_clicks + ?", count)).Error
}

// Delete removes code related url record from database
func (r *urlRepo) Delete(code string) error {
	if err := r.db.Where("code= ?", code).Delete(&model.URL{}).Error; err != nil {
		return err
	}
	return nil
}

// DeleteUserAllUrls removes all userId related url records from database
func (r *urlRepo) DeleteUserAllUrls(userID uint) ([]model.URL, error) {
	var urls []model.URL

	if err := r.db.Where("user_id = ?", userID).Find(&urls).Error; err != nil {
		return nil, err
	}

	if err := r.db.Where("user_id = ?", userID).Delete(&model.URL{}).Error; err != nil {
		return nil, err
	}

	return urls, nil
}
