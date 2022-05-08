package es_model

import "time"

var rentalInfoMapping = `
{
  "mappings": {
    "properties": {
      "id": {
        "type": "long"
      },
      "uid": {
        "type": "keyword"
      },
      "title": {
        "type": "text",
        "analyzer": "ik_max_word",
        "search_analyzer": "ik_smart"
      },
      "cover": {
        "type": "keyword"
      },
      "pics": {
        "type": "text"
      },
      "area": {
        "type": "long"
      },
      "price": {
        "type": "long"
      },
      "rent_avail_time": {
        "type": "keyword"
      },
      "rent_term_from": {
        "type": "integer"
      },
      "rent_term_to": {
        "type": "integer"
      },
      "province": {
        "type": "keyword"
      },
      "city": {
        "type": "keyword"
      },
      "location": {
        "type": "text"
      },
      "desc": {
        "type": "text",
        "analyzer": "ik_max_word",
        "search_analyzer": "ik_smart"
      },
      "tags": {
        "type": "text",
        "analyzer": "ik_max_word",
        "search_analyzer": "ik_smart"
      },
      "house_type": {
        "type": "keyword"
      },
      "room_type": {
        "type": "keyword"
      },
      "furniture": {
        "type": "text",
        "analyzer": "ik_max_word",
        "search_analyzer": "ik_smart"
      },
      "type": {
        "type": "keyword"
      },
      "rooms": {
        "type": "nested",
		"properties": {
			"RoomName": {
				"type": "text",
				"analyzer": "ik_max_word",
        		"search_analyzer": "ik_smart"
			},
			"Area": {
				"type": "long"
			},
			"Price": {
				"type": "long"	
			}
		}
      },
      "subs_num": {
        "type": "long"
      },
      "insert_time": {
        "type": "date"
      },
      "update_time": {
        "type": "date"
      },
      "is_delete": {
        "type": "long"
      }
    }
  }
}
`

//  租房信息表
type RentalInfo struct {
	RequestId     string     `gorm:"-" json:"-"`
	ID            uint       `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`                                          // 主键 ID
	Uid           string     `gorm:"column:uid;type:varchar(256);NOT NULL" json:"uid"`                                        // 发布消息用户
	Title         string     `gorm:"column:title;type:varchar(128);NOT NULL" json:"title"`                                    // 标题
	ObfsTitle     string     `json:"obfs_title"`                                                                              // 混淆标题
	Cover         string     `gorm:"column:cover;type:varchar(1024);NOT NULL" json:"cover"`                                   // 封面 (COS 链接)
	Pics          []string   `gorm:"column:pics;type:text;NOT NULL" json:"pics"`                                              // 房子照片 (COS 链接)
	Area          float64    `gorm:"column:area;type:float;NOT NULL" json:"area"`                                             // 面积
	Price         []float64  `gorm:"column:price;type:float;NOT NULL" json:"price"`                                           // 价格
	RentAvailTime string     `gorm:"column:rent_avail_time;type:varchar(64);NOT NULL" json:"rent_avail_time"`                 // 入住时间
	RentTermFrom  int        `gorm:"column:rent_term_from;NOT NULL" json:"rent_term_from"`                                    // 租房周期(最低 n 个月)
	RentTermTo    int        `gorm:"column:rent_term_to;NOT NULL" json:"rent_term_to"`                                        // 租房周期(最高 n 个月)
	Province      string     `gorm:"column:province;type:varchar(256);NOT NULL" json:"province"`                              // 省份
	City          string     `gorm:"column:city;type:varchar(256);NOT NULL" json:"city"`                                      // 城市
	Location      string     `gorm:"column:location;type:varchar(256);NOT NULL" json:"location"`                              //  位置
	Desc          string     `gorm:"column:desc;type:text;NOT NULL" json:"desc"`                                              //  描述
	ObfsDesc      string     `json:"obfs_desc"`                                                                               // 混淆描述
	Tags          []string   `gorm:"column:tags;type:text;NOT NULL" json:"tags"`                                              //  标签 (json 数组)
	HouseType     string     `gorm:"column:house_type;type:varchar(256);NOT NULL" json:"house_type"`                          //  户型
	RoomType      []string   `gorm:"column:room_type;type:varchar(1024);NOT NULL" json:"room_type"`                           //  房型 (图片链接)
	Furniture     []string   `gorm:"column:furniture;type:text;NOT NULL" json:"furniture"`                                    //  家具 (json 数组)
	Type          string     `gorm:"column:type;type:int(11);NOT NULL" json:"type"`                                           //  0 - 整租 1 - 合租
	Rooms         []RoomInfo `gorm:"column:rooms;type:text;NOT NULL" json:"rooms"`                                            //  如果为合租，出租的房间
	SubsNum       int        `gorm:"column:subs_num;type:int(11);default:0;NOT NULL" json:"subs_num"`                         // 关注数
	InsertTime    time.Time  `gorm:"column:insert_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"insert_time"` // 插入时间
	UpdateTime    time.Time  `gorm:"column:update_time;type:timestamp;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_time"` // 更新时间
	IsDelete      int        `gorm:"column:is_delete;type:tinyint(1);default:0;NOT NULL" json:"is_delete"`                    // 是否删除 0 - 未删除 1 - 删除
}

type RoomInfo struct {
	RoomName string  `json:"RoomName,omitempty" validate:"required"`
	Area     float64 `json:"Area,omitempty" validate:"required"`
	Price    float64 `json:"Price,omitempty" validate:"required"`
}

func (r RentalInfo) Index() string {
	return "rental_info"
}

func (r RentalInfo) Mapping() string {
	return rentalInfoMapping
}
