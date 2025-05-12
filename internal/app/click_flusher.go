package app

import "github.com/CemAkan/url-shortener/internal/repository"

type ClickFlusherService struct {
	repo repository.URLRepository
}

func NewClickFlusherService(repo repository.URLRepository) *ClickFlusherService {
	return &ClickFlusherService{
		repo: repo,
	}
}
