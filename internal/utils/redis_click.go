package utils

import (
	"context"
	"fmt"
	"github.com/CemAkan/url-shortener/config"
	"time"
)

// TrackClick increments short url redis click counter
func TrackClick(ctx context.Context, code string) {
	key := fmt.Sprintf("clicks:%s", code)

	pipe := config.Redis.TxPipeline()

	// increment
	pipe.Incr(ctx, key)

	pipe.Expire(ctx, key, time.Hour*24)
	_, _ = pipe.Exec(ctx)
}
