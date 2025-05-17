package repository

import (
	"github.com/CemAkan/url-shortener/internal/domain/entity"
	"github.com/CemAkan/url-shortener/internal/infrastructure/db"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entity.User) error
	FindByID(id uint) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	Update(user *entity.User) error
	ListAllUsers() ([]entity.User, error)
	GetByID(id uint) (*entity.User, error)
	Delete(id uint) error
	SetTrueMailConfirmationStatus(id uint) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepo{
		db: db.DB,
	}
}

// Create inserts new user
func (r *userRepo) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

// FindByID retrieves user by ID
func (r *userRepo) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	return &user, err
}

// FindByEmail retrieves user by username
func (r *userRepo) FindByEmail(email string) (*entity.User, error) {
	var user entity.User

	err := r.db.Where("email = ?", email).First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update modifies to existing user
func (r *userRepo) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

// ListAllUsers retrieves all user records
func (r *userRepo) ListAllUsers() ([]entity.User, error) {
	var users []entity.User

	err := r.db.Find(&users).Error

	return users, err
}

// GetByID fund user record with id parameter
func (r *userRepo) GetByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	return &user, err
}

// Delete removes user record from database
func (r *userRepo) Delete(id uint) error {
	if err := r.db.Where("id = ?", id).Delete(&entity.User{}).Error; err != nil {
		return err
	}
	return nil
}

// SetTrueMailConfirmationStatus set true is_mail_confirmed field in selected user record with userID
func (r *userRepo) SetTrueMailConfirmationStatus(id uint) error {
	if err := r.db.Model(&entity.User{}).Where("id=?", id).Updates(entity.User{IsMailConfirmed: true}).Error; err != nil {
		return err
	}
	return nil
}
