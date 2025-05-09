package repository

import (
	"github.com/CemAkan/url-shortener/config"
	"github.com/CemAkan/url-shortener/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id uint) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	ListAllUsers() ([]domain.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepo{
		db: config.DB
}
