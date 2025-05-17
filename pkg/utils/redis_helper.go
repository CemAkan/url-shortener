package utils

import (
	"context"
	"fmt"
	"github.com/CemAkan/url-shortener/internal/infrastructure/cache"
	"github.com/CemAkan/url-shortener/pkg/logger"
)

func DeleteURLCache(code string) {
	keys := []string{
		fmt.Sprintf("clicks:%s", code),
		fmt.Sprintf("hotlink:%s", code),
	}

	for _, key := range keys {
		if err := cache.Redis.Del(context.Background(), key).Err(); err != nil {
			logger.Log.Warnf("Failed to delete Redis key %s: %v", key, err)
		}
	}
}
