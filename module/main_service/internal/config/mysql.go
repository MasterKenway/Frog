package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	mysqlCli *gorm.DB
)

func GetMysqlCli() *gorm.DB {
	if mysqlCli != nil {
		return mysqlCli
	}

	db, err := gorm.Open(mysql.Open("user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	mysqlCli = db
	return db
}
