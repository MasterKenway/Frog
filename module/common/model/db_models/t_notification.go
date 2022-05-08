package db_models

import (
	"database/sql"
	"time"
)

// 用户通知表
type Notification struct {
	Id         uint          `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT" json:"id"`                    // 主键 ID
	Type       int           `gorm:"column:type;type:int(11);NOT NULL" json:"type"`                                           // 通知类型
	Uid        string        `gorm:"column:uid;type:int(11);NOT NULL" json:"uid"`                                             // 通知用户
	CommentId  sql.NullInt32 `gorm:"column:comment_id;type:int(11) unsigned" json:"comment_id"`                               // 评论
	RentalId   sql.NullInt32 `gorm:"column:rental_id;type:int(11) unsigned" json:"rental_id"`                                 // 租房信息 ID
	Content    string        `gorm:"column:content;type:varchar(128);NOT NULL" json:"content"`                                // 通知内容
	IsRead     int           `gorm:"column:is_read;type:tinyint(1);default:0;NOT NULL" json:"is_read"`                        // 已读
	InsertTime time.Time     `gorm:"column:insert_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"insert_time"` // 插入时间
	UpdateTime time.Time     `gorm:"column:update_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_time"` // 更新时间
	IsDelete   int           `gorm:"column:is_delete;type:tinyint(1);default:0;NOT NULL" json:"is_delete"`                    // 是否删除 0 - 未删除 1 - 删除
}

func (m *Notification) TableName() string {
	return "t_notification"
}
