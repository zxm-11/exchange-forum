package config

import (
	"context"
	"exchangeapp/global"
	"log"

	"github.com/go-redis/redis/v8"
)

// Redis连接配置
func InitRedis() {

	RedisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "",
	})

	if _, err := RedisClient.Ping(context.Background()).Result(); err != nil { //测试Redis连接是否成功
		log.Fatalf("Failed to connect to Redis,got error: %v", err)
	}

	global.RedisDB = RedisClient
}
