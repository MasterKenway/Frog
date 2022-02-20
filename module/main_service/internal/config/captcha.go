package config

import (
	"context"
	"encoding/json"
	"graduation-project/module/common/constant"
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

	resp, err := GetEtcdCli().Get(context.Background(), constant.EtcdKeyCaptchaConfig)
	if err != nil {
		panic(err)
	}

	if len(resp.Kvs) <= 0 {
		panic(err)
	}

	captchaConfig = &CaptchaConfig{}
	err = json.Unmarshal(resp.Kvs[0].Value, captchaConfig)
	if err != nil {
		panic(err)
	}

	return captchaConfig
}
