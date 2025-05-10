package app

import (
	"errors"
	"github.com/CemAkan/url-shortener/internal/domain"
	"github.com/CemAkan/url-shortener/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(username, password string) (*domain.User, error)
	Login(username, password string) (*domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		repo: userRepo,
	}
}

// Register checks username existence and save new user record to db with hashed password
func (s *userService) Register(username, password string) (*domain.User, error) {
	// username existence checking
	existing, err := s.repo.FindByUsername(username)

	if existing != nil && err == nil {
		return nil, errors.New("username already taken")
	}

	//password hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, errors.New("password hashing failure")
	}
	user := &domain.User{
		Username: username,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, errors.New("user create failure")
	}

	return user, nil

}

// Login checks username existence and its related password's correctness
func (s *userService) Login(username, password string) (*domain.User, error) {

	user, err := s.repo.FindByUsername(username)

	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}
