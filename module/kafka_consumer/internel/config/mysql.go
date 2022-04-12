package config

import (
	"encoding/json"
	"fmt"

	"frog/module/common/config"
	"frog/module/common/constant"

	perrors "github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mysqlCli *gorm.DB
)

func GetReadOnlyMysqlCli() *gorm.DB {
	if mysqlCli != nil {
		return mysqlCli
	}

	mysqlEtcdConfigBytes, err := GetConfig(constant.EtcdKeyMysqlConfig)
	if err != nil {
		panic(err)
	}

	var mysqlConfMap config.MysqlConfig
	err = json.Unmarshal(mysqlEtcdConfigBytes, &mysqlConfMap)
	if err != nil {
		panic(err)
	}

	mysqlConf, ok := mysqlConfMap.ConfigMaps[constant.MysqlUserReadOnly]
	if !ok {
		panic(perrors.New("mysql config key not exists"))
	}

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/dbname?charset=utf8mb4&parseTime=True&loc=Local", mysqlConf.User, mysqlConf.Password, mysqlConf.Endpoint, mysqlConf.Port)), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	mysqlCli = db
	return db
}

func GetMysqlCli() *gorm.DB {
	if mysqlCli != nil {
		return mysqlCli
	}

	mysqlEtcdConfigBytes, err := GetConfig(constant.EtcdKeyMysqlConfig)
	if err != nil {
		panic(err)
	}

	var mysqlConfMap config.MysqlConfig
	err = json.Unmarshal(mysqlEtcdConfigBytes, &mysqlConfMap)
	if err != nil {
		panic(err)
	}

	mysqlConf, ok := mysqlConfMap.ConfigMaps[constant.MysqlUserAll]
	if !ok {
		panic(perrors.New("mysql config key not exists"))
	}

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/dbname?charset=utf8mb4&parseTime=True&loc=Local", mysqlConf.User, mysqlConf.Password, mysqlConf.Endpoint, mysqlConf.Port)), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	mysqlCli = db
	return db
}
