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

	var mysqlConf *config.EmailConfig
	err = json.Unmarshal(configBytes, mysqlConf)
	if err != nil {
		panic(err)
	}

	return emailConf
}
