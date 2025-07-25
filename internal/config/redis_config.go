package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var RedisCtx = context.Background()

func RedisConfig() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     AppConfig.RedisURL,
		Password: "",
		DB:       0,
	})

	if _, err := RedisClient.Ping(RedisCtx).Result(); err != nil {
		log.Fatal("Redis connect error:", err)
	}
}
