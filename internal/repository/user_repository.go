package repository

import "github.com/CemAkan/url-shortener/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id uint) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	ListAllUsers() ([]domain.User, error)
}
