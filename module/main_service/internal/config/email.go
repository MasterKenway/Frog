package config

import (
	"encoding/json"

	"frog/module/common/config"
	"frog/module/common/constant"
)

var (
	emailConf *config.EmailConfig
)

func GetEmailConfig() *config.EmailConfig {
	if emailConf != nil {
		return emailConf
	}

	configBytes, err := GetConfig(constant.EtcdKeyEmailConfig)
	if err != nil {
		panic(err)
	}

	conf := &config.EmailConfig{}
	err = json.Unmarshal(configBytes, &conf)
	if err != nil {
		panic(err)
	}
	emailConf = conf

	return emailConf
}
