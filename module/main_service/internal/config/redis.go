package config

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
)

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}
}

func GetRedisCli() *redis.Client {
	return redisClient
}
