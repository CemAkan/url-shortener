package app

import "github.com/CemAkan/url-shortener/internal/domain"

type UserService interface {
	Register(username, password string) (*domain.User, error)
	Login(username, password string) (*domain.User, error)
}
