package app

import (
	"context"
	"errors"
	"github.com/CemAkan/url-shortener/config"
	"github.com/CemAkan/url-shortener/internal/domain"
	"github.com/CemAkan/url-shortener/internal/repository"
	"github.com/CemAkan/url-shortener/internal/utils"
)

type URLService interface {
	Shorten(originalURL string, userID uint, customCode *string) (*domain.URL, error)
	GetByCode(code string) (*domain.URL, error)
	GetUserURLs(userID uint) ([]domain.URL, error)
	GetSingleUrlRecord(code string, userID uint) (*domain.URL, int, error)
}

type urlService struct {
	repo repository.URLRepository
}

func NewURLService(urlRepo repository.URLRepository) URLService {
	return &urlService{
		repo: urlRepo,
	}
}

// Shorten redeclare url address
func (s *urlService) Shorten(originalURL string, userID uint, customCode *string) (*domain.URL, error) {
	var code string

	//custom uniqueness check
	if customCode != nil && *customCode != "" {
		existing, _ := s.repo.FindByCode(*customCode)
		if existing != nil && existing.ID != 0 {
			return nil, errors.New("custom code already taken")
		}
		code = *customCode
	} else {
		// unique code generator
		for {
			code = utils.GenerateCode(7)
			existing, _ := s.repo.FindByCode(code)
			if existing == nil || existing.ID == 0 {
				break
			}
		}
	}

	//create and assert url
	url := &domain.URL{
		Code:        code,
		OriginalURL: originalURL,
		UserID:      userID,
	}

	// save to db new url parameter
	if err := s.repo.Create(url); err != nil {
		return nil, err
	}

	return url, nil

}

// GetByCode finds url record with code parameter
func (s *urlService) GetByCode(code string) (*domain.URL, error) {
	return s.repo.FindByCode(code)
}

// GetUserURLs finds all userID related url records
func (s *urlService) GetUserURLs(userID uint) ([]domain.URL, error) {
	return s.repo.FindByUserID(userID)

}

func (s *urlService) GetSingleUrlRecord(code string, userID uint) (*domain.URL, int, error) {
	url, err := s.repo.FindByCode(code)
	if err != nil || url == nil {
		return nil, 0, err
	}
	if url.UserID != userID {
		return nil, 0, errors.New("unauthorized access")
	}

	// getting daily click rate from redis
	clickKey := "clicks:" + code
	dailyClicks, _ := config.Redis.Get(context.Background(), clickKey).Int()

	return url, dailyClicks, nil
}
