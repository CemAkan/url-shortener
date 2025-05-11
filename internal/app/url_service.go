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
