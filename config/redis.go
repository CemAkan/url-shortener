package config

import (
	redis2 "github.com/redis/go-redis/v9"
)

var Redis *redis2.Client

func InitRedis() {
	addr := GetEnv("REDIS_ADDR", "")
	
}
