package config

import (
	"encoding/json"

	"frog/module/common/config"
	"frog/module/common/constant"
)

var (
	rateLimitConf *config.RateLimitConfig
)

func GetRateLimitConfig() *config.RateLimitConfig {
	if rateLimitConf != nil {
		return rateLimitConf
	}

	conf := config.RateLimitConfig{}
	confBytes, err := GetConfig(constant.EtcdKeyRateLimitConfig)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(confBytes, &conf)
	if err != nil {
		panic(err)
	}

	rateLimitConf = &conf
	return rateLimitConf
}
