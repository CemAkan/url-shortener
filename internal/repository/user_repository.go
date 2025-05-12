package repository

import (
	"github.com/CemAkan/url-shortener/internal/domain"
	"github.com/CemAkan/url-shortener/pkg/infrastructure"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id uint) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	Update(user *domain.User) error
	ListAllUsers() ([]domain.User, error)
	GetByID(id uint) (*domain.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepo{
		db: infrastructure.DB,
	}
}

// Create inserts new user
func (r *userRepo) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

// FindByID retrieves user by ID
func (r *userRepo) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	return &user, err
}

// FindByUsername retrieves user by username
func (r *userRepo) FindByUsername(username string) (*domain.User, error) {
	var user domain.User //replacing

	err := r.db.Where("username = ?", username).First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update modifies to existing user
func (r *userRepo) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

// ListAllUsers retrieves all user records
func (r *userRepo) ListAllUsers() ([]domain.User, error) {
	var users []domain.User

	err := r.db.Find(&users).Error

	return users, err
}

// GetByID fund user record with id parameter
func (r *userRepo) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	return &user, err
}
