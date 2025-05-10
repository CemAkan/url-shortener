package app

import (
	"errors"
	"github.com/CemAkan/url-shortener/internal/domain"
	"github.com/CemAkan/url-shortener/internal/repository"
	"github.com/CemAkan/url-shortener/utils"
)

type URLService interface {
	Shorten(originalURL string, userID uint, customCode *string) (*domain.URL, error)
	GetByCode(code string) (*domain.URL, error)
	GetByUserID(userID uint) ([]domain.URL, error)
}

type urlService struct {
	repo repository.URLRepository
}

func NewURLService(urlRepo repository.URLRepository) URLService {
	return &userService{
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
