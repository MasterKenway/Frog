package db_models

import (
	"time"
)

// 日志表
type Log struct {
	Id         uint      `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT" json:"id"`                    // 主键 ID
	Time       time.Time `gorm:"column:time;type:timestamp;NOT NULL" json:"time"`                                         // log 记录时间
	Level      string    `gorm:"column:level;type:varchar(16);NOT NULL" json:"level"`                                     // log 等级
	Caller     string    `gorm:"column:caller;type:varchar(256);NOT NULL" json:"caller"`                                  // 调用函数
	RequestId  string    `gorm:"column:request_id;type:varchar(256);NOT NULL" json:"request_id"`                          // 请求 ID
	Message    string    `gorm:"column:message;type:text;NOT NULL" json:"message"`                                        // 消息
	InsertTime time.Time `gorm:"column:insert_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"insert_time"` // 插入时间
	UpdateTime time.Time `gorm:"column:update_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_time"` // 更新时间
	IsDelete   int       `gorm:"column:is_delete;type:tinyint(1);default:0;NOT NULL" json:"is_delete"`                    // 是否删除 0 - 未删除 1 - 删除
}

func (m *Log) TableName() string {
	return "l_log"
}
