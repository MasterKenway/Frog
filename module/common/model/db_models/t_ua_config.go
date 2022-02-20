package db_models

import (
	"time"
)

type UaConfig struct {
	ID         uint      `gorm:"column:id;primary_key;AUTO_INCREMENT"`                  // 主键 ID
	Ua         string    `gorm:"column:ua;NOT NULL"`                                    // UA 黑名单配置
	InsertTime time.Time `gorm:"column:insert_time;default:CURRENT_TIMESTAMP;NOT NULL"` // 插入时间
	UpdateTime time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;NOT NULL"` // 更新时间
	IsDelete   int       `gorm:"column:is_delete;default:0;NOT NULL"`                   // 是否删除 0 - 未删除 1 - 删除
}

func (m *UaConfig) TableName() string {
	return "t_ua_config"
}
