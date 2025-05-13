package app

import (
	"errors"
	"github.com/CemAkan/url-shortener/internal/domain/model"
	"github.com/CemAkan/url-shortener/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"strings"
	"unicode"
)

type UserService interface {
	Register(email, password, name, surname string) (*model.User, error)
	Login(email, password string) (*model.User, error)
	GetByID(id uint) (*model.User, error)
	DeleteUser(id uint) error
	ListAllUsers() ([]model.User, error)
	SetTrueEmailConfirmation(id uint) error
	PasswordUpdate(userID uint, newPassword string) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		repo: userRepo,
	}
}

// Register checks email existence and save new user record to db with hashed password
func (s *userService) Register(email, password, name, surname string) (*model.User, error) {
	// email existence checking
	existing, err := s.repo.FindByEmail(email)

	if existing != nil && err == nil {
		return nil, errors.New("email already registered")
	}

	//password hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, errors.New("password hashing failure")
	}

	formatedName, err := s.format(name)
	if err != nil {
		return nil, errors.New("name required")
	}

	formatedSurname, err := s.format(surname)
	if err != nil {
		return nil, errors.New("surname required")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errors.New("invalid email address")
	}

	user := &model.User{
		Name:     formatedName,
		Surname:  formatedSurname,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, errors.New("user create failure")
	}

	return user, nil

}

// Login checks email existence and its related password's correctness
func (s *userService) Login(email, password string) (*model.User, error) {

	user, err := s.repo.FindByEmail(email)

	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

// GetByID return user with its id
func (s *userService) GetByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

// DeleteUser deletes user record
func (s *userService) DeleteUser(id uint) error {
	exist, _ := s.GetByID(id)
	if exist == nil {
		return errors.New("user not exist")
	}

	err := s.repo.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

// ListAllUsers gets all user records
func (s *userService) ListAllUsers() ([]model.User, error) {
	return s.repo.ListAllUsers()
}

// Format cleans spaces, lower-cases everything, then capitalises only the first rune.
func (s *userService) format(word string) (string, error) {
	// remove spaces outside
	w := strings.ToLower(strings.TrimSpace(word))

	// remove spaces inside
	w = strings.ReplaceAll(w, " ", "")

	//if it is empty return it
	if w == "" {
		return "", errors.New("empty input")
	}

	//make capital first char
	runes := []rune(w)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes), nil
}

// SetTrueEmailConfirmation sets true is_email_confirmed field
func (s *userService) SetTrueEmailConfirmation(id uint) error {
	return s.repo.SetTrueMailConfirmationStatus(id)
}

// PasswordUpdate change password with new one
func (s *userService) PasswordUpdate(userID uint, newPassword string) error {
	user, err := s.GetByID(userID)

	if err != nil {
		return errors.New("user not found")
	}

	user.Password = newPassword

	return s.repo.Update(user)
}
