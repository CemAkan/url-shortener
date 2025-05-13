package repository

import (
	"github.com/CemAkan/url-shortener/internal/domain/model"
	"github.com/CemAkan/url-shortener/pkg/infrastructure"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByID(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	ListAllUsers() ([]model.User, error)
	GetByID(id uint) (*model.User, error)
	Delete(id uint) error
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
func (r *userRepo) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// FindByID retrieves user by ID
func (r *userRepo) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

// FindByEmail retrieves user by username
func (r *userRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User

	err := r.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update modifies to existing user
func (r *userRepo) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// ListAllUsers retrieves all user records
func (r *userRepo) ListAllUsers() ([]model.User, error) {
	var users []model.User

	err := r.db.Find(&users).Error

	return users, err
}

// GetByID fund user record with id parameter
func (r *userRepo) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

// Delete removes user record from database
func (r *userRepo) Delete(id uint) error {
	if err := r.db.Where("id = ?", id).Delete(&model.User{}).Error; err != nil {
		return err
	}
	return nil
}
