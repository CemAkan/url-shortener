package service

import (
	"context"
	"errors"
	"github.com/CemAkan/url-shortener/config"
	"github.com/CemAkan/url-shortener/internal/domain/entity"
	"github.com/CemAkan/url-shortener/internal/infrastructure/cache"
	"github.com/CemAkan/url-shortener/internal/metrics"
	"github.com/CemAkan/url-shortener/internal/repository"
	"github.com/CemAkan/url-shortener/pkg/logger"
	"github.com/CemAkan/url-shortener/pkg/utils"
	"strconv"
	"strings"
	"time"
)

type URLService interface {
	Shorten(originalURL string, userID uint, customCode *string) (*entity.URL, error)
	GetUserURLs(userID uint) ([]entity.URL, error)
	GetSingleUrlRecord(code string) (*entity.URL, int, error)
	ResolveRedirect(ctx context.Context, code string) (string, error)
	UpdateUserURL(userID uint, oldCode string, newOriginalURL, newCode *string) error
	DeleteUserURL(userID uint, code string) error
	DeleteUserAllURLs(userID uint) error
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
func (s *urlService) Shorten(originalURL string, userID uint, customCode *string) (*entity.URL, error) {
	var code string

	if customCode != nil && *customCode != "" {
		isTaken := s.isCodeTaken(*customCode)

		if isTaken {
			return nil, errors.New("custom code already taken")
		}
		code = *customCode
	} else {
		code = s.generateUniqueCode()
	}

	// if it is api/* or api, generate it
	if strings.Contains(code, "api/") || code == "api" {
		code = s.generateUniqueCode()
	}

	url := &entity.URL{
		Code:        code,
		OriginalURL: originalURL,
		UserID:      userID,
	}

	if err := s.repo.Create(url); err != nil {
		return nil, err
	}

	return url, nil
}

func (s *urlService) isCodeTaken(code string) bool {
	existing, _ := s.repo.FindByCode(code)

	if existing != nil && existing.ID != 0 {
		return true
	}

	return false
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

// GetUserURLs finds all userID related url records
func (s *urlService) GetUserURLs(userID uint) ([]entity.URL, error) {
	return s.repo.FindByUserID(userID)

}

// GetSingleUrlRecord find url record with its daily click rate
func (s *urlService) GetSingleUrlRecord(code string) (*entity.URL, int, error) {
	url, err := s.repo.FindByCode(code)
	if err != nil || url == nil {
		return nil, 0, err
	}

	// getting daily click rate from redis
	clickKey := "clicks:" + code
	dailyClicks, _ := cache.Redis.Get(context.Background(), clickKey).Int()

	return url, dailyClicks, nil
}

// ResolveRedirect translates given short code to original code with cache-db mechanism
func (s *urlService) ResolveRedirect(ctx context.Context, code string) (string, error) {

	// prometheus latency
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RedirectLatency.Observe(duration)
	}()

	// prometheus click counter
	metrics.ClickCounter.WithLabelValues(code).Inc()

	//look at redis to find cache record
	cacheKey := "code_cache:" + code
	if originalURL, err := cache.Redis.Get(ctx, cacheKey).Result(); err == nil && originalURL != "" {
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
		if err := cache.Redis.Set(ctx, cacheKey, url.OriginalURL, 24*time.Hour).Err(); err != nil {
			logger.Log.Printf("Redis cache save error: %v", err.Error())
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
		isTaken := s.isCodeTaken(*newCode)

		if isTaken {
			return errors.New("new custom code already taken")
		}
		url.Code = *newCode
	}

	//original url update check
	if newOriginalURL != nil && *newOriginalURL != "" && *newOriginalURL != url.OriginalURL {
		url.OriginalURL = *newOriginalURL
	}

	// redis cache clean
	utils.DeleteURLCache(oldCode)

	return s.repo.Update(url)
}

func (s *urlService) DeleteUserURL(userID uint, code string) error {
	url, err := s.repo.FindByCode(code)
	if err != nil {
		return errors.New("url not found")
	}

	if url.UserID != userID {
		return errors.New("unauthorized")
	}

	utils.DeleteURLCache(code)

	return s.repo.Delete(code)
}

// DeleteUserAllURLs removes user relational urls and their redis data
func (s *urlService) DeleteUserAllURLs(userID uint) error {

	urls, err := s.repo.DeleteUserAllUrls(userID)

	if err != nil {
		return err
	}

	// Redis key cleanup
	for _, url := range urls {
		utils.DeleteURLCache(url.Code)
	}

	return nil

}
