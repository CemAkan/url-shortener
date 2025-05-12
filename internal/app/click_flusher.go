package app

import (
	"context"
	"github.com/CemAkan/url-shortener/internal/repository"
	"github.com/CemAkan/url-shortener/internal/utils"
	"github.com/CemAkan/url-shortener/pkg/infrastructure"
)

type ClickFlusherService struct {
	repo repository.URLRepository
}

func NewClickFlusherService(repo repository.URLRepository) *ClickFlusherService {
	return &ClickFlusherService{
		repo: repo,
	}
}

func (s *ClickFlusherService) FlushClicks() {
	ctx := context.Background()
	keys, err := utils.GetAllClickKeys(ctx)

	if err != nil {
		infrastructure.Log.WithError(err).Error("Failed to get click keys from Redis")
		return
	}

	for _, key := range keys {
		code := key[len("clicks:"):] //get after 7 char

		count, err := utils.GetDailyClickCount(ctx, code)
		if err != nil {
			infrastructure.Log.WithError(err).Warnf("Failed to get count for %s", key)
			continue
		}

		if err := s.repo.AddToTotalClicks(code, count); err != nil {
			infrastructure.Log.WithError(err).Errorf("Failed to update DB clicks for %s", code)
			continue
		}

		if err := utils.DeleteClickKey(ctx, code); err != nil {
			infrastructure.Log.WithError(err).Warnf("Failed to delete Redis key %s", key)
			continue
		}

		infrastructure.Log.Infof("Flushed %d clicks for %s", count, code)
	}
}
