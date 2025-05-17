package cache

import (
	"context"
	"github.com/CemAkan/url-shortener/config"
	"github.com/CemAkan/url-shortener/pkg/infrastructure/logger"
	"github.com/redis/go-redis/v9"
	"time"
)

var Redis *redis.Client

func InitRedis() {
	addr := config.GetEnv("REDIS_ADDR", "")
	password := config.GetEnv("REDIS_PASSWORD", "")

	Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	err := Redis.Ping(ctx).Err()

	if err != nil {
		logger.Log.Fatalf("Failed to connect to redis: %v", err)
	}

	logger.Log.Info("Redis connection established successfully")

}
