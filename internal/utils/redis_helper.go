package utils

import (
	"context"
	"fmt"
	"github.com/CemAkan/url-shortener/pkg/infrastructure"
)

func DeleteURLCache(code string) {
	keys := []string{
		fmt.Sprintf("clicks:%s", code),
		fmt.Sprintf("hotlink:%s", code),
	}

	for _, key := range keys {
		if err := infrastructure.Redis.Del(context.Background(), key).Err(); err != nil {
			infrastructure.Log.Warnf("Failed to delete Redis key %s: %v", key, err)
		}
	}
}
