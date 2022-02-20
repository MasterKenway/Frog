package config

import (
	"encoding/json"
	"fmt"
	"graduation-project/module/common/constant"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"graduation-project/module/common/config"
)

var (
	mysqlCli *gorm.DB
)

func GetMysqlCli() *gorm.DB {
	if mysqlCli != nil {
		return mysqlCli
	}

	mysqlEtcdConfigBytes, err := GetConfig(constant.EtcdKeyMysqlConfig)
	if err != nil {
		panic(err)
	}

	var mysqlConf config.MysqlConfig
	err = json.Unmarshal(mysqlEtcdConfigBytes, &mysqlConf)
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/dbname?charset=utf8mb4&parseTime=True&loc=Local", mysqlConf.User, mysqlConf.Password, mysqlConf.Endpoint, mysqlConf.Port)), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	mysqlCli = db
	return db
}
