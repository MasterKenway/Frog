package config

import (
	"encoding/json"

	"frog/module/common/constant"
)

var (
	captchaConfig *CaptchaConfig
)

type CaptchaConfig struct {
	AppID        string `json:"app_id,omitempty"`
	AppSecretKey string `json:"app_secret_key,omitempty"`
	ApiSecretID  string `json:"api_secret_id,omitempty"`
	ApiSecretKey string `json:"api_secret_key,omitempty"`
}

func GetCaptchaConfig() *CaptchaConfig {
	if captchaConfig != nil {
		return captchaConfig
	}

	confBytes, err := GetConfig(constant.EtcdKeyCaptchaConfig)
	if err != nil {
		panic(err)
	}

	captchaConfig = &CaptchaConfig{}
	err = json.Unmarshal(confBytes, captchaConfig)
	if err != nil {
		panic(err)
	}

	return captchaConfig
}
