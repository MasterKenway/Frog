package service

import (
	"time"

	"frog/module/kafka_consumer/internel/config"
	"frog/module/kafka_consumer/internel/log"

	"frog/module/common/constant"
	"frog/module/common/model/db_models"
)

func DeleteLog() {
	date := time.Now().Add(-1 * time.Duration(30) * constant.OneDay)
	err := config.GetMysqlCli().Where("insert_time < ?", date).Delete(&db_models.Log{}).Error
	if err != nil {
		log.Errorf("failed to delete logs, %s", err.Error())
	}
}
