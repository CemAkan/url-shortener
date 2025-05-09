package config

import (
	redis2 "github.com/redis/go-redis/v9"
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
}
