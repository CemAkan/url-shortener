package app

import (
	"context"
	"errors"
	"github.com/CemAkan/url-shortener/config"
	"github.com/CemAkan/url-shortener/internal/domain"
	"github.com/CemAkan/url-shortener/internal/repository"
	"github.com/CemAkan/url-shortener/internal/utils"
	"strconv"
	"time"
)

type URLService interface {
	Shorten(originalURL string, userID uint, customCode *string) (*domain.URL, error)
	GetByCode(code string) (*domain.URL, error)
	GetUserURLs(userID uint) ([]domain.URL, error)
	GetSingleUrlRecord(code string, userID uint) (*domain.URL, int, error)
	ResolveRedirect(ctx context.Context, code string) (string, error)
	UpdateUserURL(userID uint, code, newOriginalURL string) error
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

	if customCode != nil && *customCode != "" {
		isTaken, err := s.isCodeTaken(*customCode)
		if err != nil {
			return nil, err
		}
		if isTaken {
			return nil, errors.New("custom code already taken")
		}
		code = *customCode
	} else {
		code = s.generateUniqueCode()
	}

	url := &domain.URL{
		Code:        code,
		OriginalURL: originalURL,
		UserID:      userID,
	}

	if err := s.repo.Create(url); err != nil {
		return nil, err
	}

	return url, nil
}

func (s *urlService) isCodeTaken(code string) (bool, error) {
	existing, err := s.repo.FindByCode(code)
	if err != nil {
		return false, err
	}
	return existing != nil && existing.ID != 0, nil
}

func (s *urlService) generateUniqueCode() string {
	for {
		code := utils.GenerateCode(7)
		existing, _ := s.repo.FindByCode(code)
		if existing == nil || existing.ID == 0 {
			return code
		}
	}
}

// GetByCode finds url record with code parameter
func (s *urlService) GetByCode(code string) (*domain.URL, error) {
	return s.repo.FindByCode(code)
}

// GetUserURLs finds all userID related url records
func (s *urlService) GetUserURLs(userID uint) ([]domain.URL, error) {
	return s.repo.FindByUserID(userID)

}

// GetSingleUrlRecord find a url record with its daily click rate
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

// ResolveRedirect translates given short code to original code with cache-db mechanism
func (s *urlService) ResolveRedirect(ctx context.Context, code string) (string, error) {

	//look at redis to find cache record
	cacheKey := "code_cache:" + code
	if originalURL, err := config.Redis.Get(ctx, cacheKey).Result(); err == nil && originalURL != "" {
		return originalURL, nil
	}

	//get daily click
	dailyClicks, _ := utils.GetDailyClickCount(ctx, code)

	//look at db to find original record
	url, err := s.repo.FindByCode(code)
	if err != nil || url == nil {
		return "", errors.New("not found")
	}

	//getting threshold from .env and transfer it to integer
	thresholdENVString := config.GetEnv("DAILY_CLICK_CACHE_THRESHOLD", "100")
	threshold, _ := strconv.Atoi(thresholdENVString)

	// db -> redis resolve redirect mechanism for hot links
	if dailyClicks >= threshold {
		if err := config.Redis.Set(ctx, cacheKey, url.OriginalURL, 24*time.Hour).Err(); err != nil {
			config.Log.Printf("Redis cache save error: %v", err.Error())
		}

	}

	return url.OriginalURL, nil

}

// UpdateUserURL updates specific code record values
func (s *urlService) UpdateUserURL(userID uint, oldCode string, newOriginalURL, newCode *string) error {
	url, err := s.repo.FindByCode(oldCode)
	if err != nil {
		return errors.New("url not found")
	}

	if url.UserID != userID {
		return errors.New("unauthorized")
	}

	// code update check
	if newCode != nil && *newCode != "" && *newCode != url.Code {
		isTaken, err := s.isCodeTaken(*newCode)
		if err != nil {
			return err
		}
		if isTaken {
			return errors.New("new custom code already taken")
		}
		url.Code = *newCode
	}

	//original url update check
	if newOriginalURL != nil && *newOriginalURL != "" && *newOriginalURL != url.OriginalURL {
		url.OriginalURL = *newOriginalURL
	}
	return s.repo.Update(url)
}
