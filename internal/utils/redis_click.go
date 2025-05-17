package utils

import (
	"context"
	"fmt"
	"github.com/CemAkan/url-shortener/pkg/infrastructure/cache"
)

// TrackClick increments short url redis click counter
func TrackClick(ctx context.Context, code string) {
	key := fmt.Sprintf("clicks:%s", code)

	pipe := cache.Redis.TxPipeline()

	// increment
	pipe.Incr(ctx, key)

	_, _ = pipe.Exec(ctx)
}

// GetDailyClickCount gets url click counter from redis
func GetDailyClickCount(ctx context.Context, code string) (int, error) {
	key := fmt.Sprintf("clicks:%s", code)
	return cache.Redis.Get(ctx, key).Int()
}

// GetAllClickKeys gets all urls click counter records from redis
func GetAllClickKeys(ctx context.Context) ([]string, error) {
	return cache.Redis.Keys(ctx, "clicks:*").Result()
}

// DeleteClickKey deletes click record
func DeleteClickKey(ctx context.Context, code string) error {
	key := fmt.Sprintf("clicks:%s", code)
	return cache.Redis.Del(ctx, key).Err()
}
