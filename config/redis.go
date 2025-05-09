package config

import (
	"context"
	redis2 "github.com/redis/go-redis/v9"
	"time"
)

var Redis *redis2.Client

func InitRedis() {
	addr := GetEnv("REDIS_ADDR", "")
	password := GetEnv("REDIS_PASSWORD", "")

	Redis = redis2.NewClient(&redis2.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	err := Redis.Ping(ctx).Err()

	if err != nil {
		Log.Fatalf("Failed to connect to redis: %v", err)
	}

	Log.Info("Redis connection established successfully")

}
