package db_models

import (
	"time"

	"frog/module/common/tools"
)

var (
	rentalInfoMap map[string]bool

	RentTypeStrMapInt = map[string]int{
		"part":  0,
		"whole": 1,
	}
)

func init() {
	rentalInfoMap = tools.GetModelCols(RentalInfo{})
}

// 租房信息表
type RentalInfo struct {
	ID            uint      `gorm:"column:id;primary_key;AUTO_INCREMENT"`                  // 主键 ID
	Uid           string    `gorm:"column:uid;NOT NULL"`                                   // 发布消息用户
	Title         string    `gorm:"column:title;NOT NULL"`                                 // 标题
	Cover         string    `gorm:"column:cover;NOT NULL"`                                 // 封面 (COS 链接)
	Pics          string    `gorm:"column:pics;NOT NULL"`                                  // 房子照片 (COS 链接)
	Area          float64   `gorm:"column:area;NOT NULL"`                                  // 面积
	Price         int       `gorm:"column:price;NOT NULL"`                                 // 价格
	RentAvailTime string    `gorm:"column:rent_avail_time;NOT NULL"`                       // 入住时间
	RentTermFrom  int       `gorm:"column:rent_term_from;NOT NULL"`                        // 租房周期(最低 n 个月)
	RentTermTo    int       `gorm:"column:rent_term_to;NOT NULL"`                          // 租房周期(最高 n 个月)
	Province      string    `gorm:"column:province;NOT NULL"`                              // 省份
	City          string    `gorm:"column:city;NOT NULL"`                                  // 城市
	Location      string    `gorm:"column:location;NOT NULL"`                              // 位置
	Desc          string    `gorm:"column:desc;NOT NULL"`                                  // 描述
	Tags          string    `gorm:"column:tags;NOT NULL"`                                  // 标签 (json 数组)
	HouseType     string    `gorm:"column:house_type;NOT NULL"`                            // 户型
	RoomType      string    `gorm:"column:room_type;NOT NULL"`                             // 房型 (图片链接)
	Furniture     string    `gorm:"column:furniture;NOT NULL"`                             // 家具 (json 数组)
	Type          int       `gorm:"column:type;NOT NULL"`                                  // 0 - 整租 1 - 合租
	Rooms         string    `gorm:"column:rooms;NOT NULL"`                                 // 如果为合租，出租的房间信息
	SubsNum       int       `gorm:"column:subs_num;default:0;NOT NULL"`                    // 关注数
	InsertTime    time.Time `gorm:"column:insert_time;default:CURRENT_TIMESTAMP;NOT NULL"` // 插入时间
	UpdateTime    time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;NOT NULL"` // 更新时间
	IsDelete      int       `gorm:"column:is_delete;default:0;NOT NULL"`                   // 是否删除 0 - 未删除 1 - 删除
}

func (m *RentalInfo) TableName() string {
	return "t_rental_info"
}

func GetRentalInfoCols() map[string]bool {
	return rentalInfoMap
}
