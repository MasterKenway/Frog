package db_models

import (
	"database/sql"
	"time"
)

// 用户表
type User struct {
	Id         uint           `gorm:"column:id;type:int(11) unsigned;primary_key;AUTO_INCREMENT" json:"id"`                    // 主键 ID
	Uid        string         `gorm:"column:uid;type:varchar(256);NOT NULL" json:"uid"`                                        // 用户 id
	Username   string         `gorm:"column:username;type:varchar(128);NOT NULL" json:"username"`                              // 用户名
	Password   string         `gorm:"column:password;type:varchar(32);NOT NULL" json:"password"`                               // 密码
	Email      string         `gorm:"column:email;type:varchar(128);NOT NULL" json:"email"`                                    // 邮箱
	Avatar     string         `gorm:"column:avatar;type:varchar(1024);NOT NULL" json:"avatar"`                                 // 头像
	LoginIps   sql.NullString `gorm:"column:login_ips;type:text" json:"login_ips"`                                             // 登录 ip （json 数组）
	IsValid    int            `gorm:"column:is_valid;type:tinyint(1);default:0;NOT NULL" json:"is_valid"`                      // 是否允许登录
	InsertTime time.Time      `gorm:"column:insert_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"insert_time"` // 插入时间
	UpdateTime time.Time      `gorm:"column:update_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_time"` // 更新时间
	IsDelete   int            `gorm:"column:is_delete;type:tinyint(1);default:0;NOT NULL" json:"is_delete"`                    // 是否删除 0 - 未删除 1 - 删除
}

func (m *User) TableName() string {
	return "t_user"
}
