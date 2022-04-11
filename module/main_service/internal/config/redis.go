package config

import (
	"context"
	"encoding/json"

	"frog/module/common/config"
	"frog/module/common/constant"

	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
)

func GetRedisCli() *redis.Client {
	if redisClient != nil {
		return redisClient
	}

	bytes, err := GetConfig(constant.EtcdKeyRedisConfig)
	if err != nil {
		panic(err)
	}

	var conf *config.RedisConfig
	err = json.Unmarshal(bytes, conf)
	if err != nil {
		panic(err)
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     conf.Endpoint,
		Password: conf.Password, // no password set
		DB:       0,             // use default DB
	})

	err = redisClient.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	return redisClient
}
