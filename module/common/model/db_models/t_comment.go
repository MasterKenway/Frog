package db_models

import (
	"time"
)

// 租房评论
type Comment struct {
	Id         uint      `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT" json:"id"`                    // 主键 ID
	RentalId   uint      `gorm:"column:rental_id;type:int(11) unsigned;NOT NULL" json:"rental_id"`                        // 住房信息 ID
	PublishUid string    `gorm:"column:publish_uid;type:int(11) unsigned;NOT NULL" json:"publish_uid"`                    // 发布人 uid
	Uid        string    `gorm:"column:uid;type:varchar(256);NOT NULL" json:"uid"`                                        // 用户 uid
	Username   string    `gorm:"column:username;type:varchar(128);NOT NULL" json:"username"`                              // 用户名称
	Avatar     string    `gorm:"column:avatar;type:varchar(1024);NOT NULL" json:"avatar"`                                 // 用户头像
	Content    string    `gorm:"column:content;type:text;NOT NULL" json:"content"`                                        // 评论内容
	InsertTime time.Time `gorm:"column:insert_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"insert_time"` // 插入时间
	UpdateTime time.Time `gorm:"column:update_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_time"` // 更新时间
	IsDelete   int       `gorm:"column:is_delete;type:tinyint(1);default:0;NOT NULL" json:"is_delete"`                    // 是否删除 0 - 未删除 1 - 删除
}

func (m *Comment) TableName() string {
	return "t_comment"
}
