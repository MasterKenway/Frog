package db_models

import (
	"time"
)

// 忽略中间件 IP 配置
type EscapeMiddlewareIps struct {
	ID         uint      `gorm:"column:id;primary_key;AUTO_INCREMENT"`                  // 主键 ID
	IP         string    `gorm:"column:ip;NOT NULL"`                                    // IP
	InsertTime time.Time `gorm:"column:insert_time;default:CURRENT_TIMESTAMP;NOT NULL"` // 插入时间
	UpdateTime time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;NOT NULL"` // 更新时间
	IsDelete   int       `gorm:"column:is_delete;default:0;NOT NULL"`                   // 是否删除 0 - 未删除 1 - 删除
}

func (m *EscapeMiddlewareIps) TableName() string {
	return "t_escape_middleware_ips"
}
