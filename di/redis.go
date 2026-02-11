package di

import (
	"digital-wallet/configs"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func SetupCache(cfg *configs.Config) *redis.Client {
	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return redisClient
}
