package db_models

import (
	"database/sql"
	"time"
)

// 用户订阅租房信息
type Subscribe struct {
	Id         uint           `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT" json:"id"`                    // 主键 ID
	Uid        string         `gorm:"column:uid;type:varchar(256);NOT NULL" json:"uid"`                                        // 用户 uid
	Type       int            `gorm:"column:type;type:int(11);NOT NULL" json:"type"`                                           // 订阅类型 0 - 租房， 1 - 标签，2 - 城市，3 - 位置
	RentalId   sql.NullInt32  `gorm:"column:rental_id;type:int(11) unsigned" json:"rental_id"`                                 // 住房信息 ID
	Tag        sql.NullString `gorm:"column:tag;type:varchar(64)" json:"tag"`                                                  // 订阅 Tag
	City       sql.NullString `gorm:"column:city;type:varchar(64)" json:"city"`                                                // 订阅城市
	Location   sql.NullString `gorm:"column:location;type:varchar(64)" json:"location"`                                        // 订阅位置
	IsValid    int            `gorm:"column:is_valid;type:tinyint(1);default:0;NOT NULL" json:"is_valid"`                      // 是否取消
	InsertTime time.Time      `gorm:"column:insert_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"insert_time"` // 插入时间
	UpdateTime time.Time      `gorm:"column:update_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_time"` // 更新时间
	IsDelete   int            `gorm:"column:is_delete;type:tinyint(1);default:0;NOT NULL" json:"is_delete"`                    // 是否删除 0 - 未删除 1 - 删除
}

func (m *Subscribe) TableName() string {
	return "t_subscribe"
}
