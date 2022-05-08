package rental_info

import (
	"encoding/json"
	"reflect"
	"strconv"

	"frog/module/common/constant"
	"frog/module/common/model/api_models"
	"frog/module/common/model/db_models"
	"frog/module/common/model/es_model"
	"frog/module/main_service/internal/config"
	"frog/module/main_service/internal/controller/auth/login"
	"frog/module/main_service/internal/log"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func GetRentalInfoCreateController() api_models.ApiInterface {
	return &RentalInfoCreateRequest{}
}

type RentalInfoCreateRequest struct {
	Title         string              `json:"Title,omitempty" validate:"required"`
	Cover         string              `json:"Cover,omitempty" validate:"required"`
	Pics          []string            `json:"Pics,omitempty" validate:"required"`
	Area          float64             `json:"Area,omitempty" validate:"required"`
	Price         int                 `json:"Price,omitempty" validate:"required"`
	RentAvailTime string              `json:"RentAvailTime,omitempty" validate:"required,oneof=ltOneMonth ltTwoMonth ltThreeMonth gtThreeMonth"`
	RentTermFrom  int                 `json:"RentTermFrom,omitempty" validate:"required,gt=0"`
	RentTermTo    int                 `json:"RentTermTo,omitempty" validate:"required,gt=1"`
	Province      string              `json:"Province,omitempty" validate:"required"`
	City          string              `json:"City,omitempty" validate:"required"`
	Location      string              `json:"Location,omitempty" validate:"required"`
	Desc          string              `json:"Desc,omitempty" validate:"required"`
	Tags          []string            `json:"Tags,omitempty"`
	HouseType     string              `json:"HouseType,omitempty" validate:"required"`
	RoomType      []string            `json:"RoomType,omitempty"`
	Furniture     []string            `json:"Furniture,omitempty"`
	Type          string              `json:"Type,omitempty" validate:"required,oneof=part whole"`
	Rooms         []es_model.RoomInfo `json:"Rooms,omitempty"`
}

func (req *RentalInfoCreateRequest) GetResult(ctx *gin.Context) (interface{}, *api_models.APIError) {
	var (
		userInfo, _ = ctx.Get(constant.CtxKeyUserInfo)
		reqId       = ctx.GetString(constant.CtxKeyRequestID)
	)

	err := saveToESAndDB(reqId, userInfo.(*login.RedisUserInfo), req)
	if err != nil {
		log.Errorf(reqId, "%s", err.Error())
		return nil, &api_models.APIError{
			Code:    constant.CodeInternalError,
			Message: "failed to create rental info",
		}
	}

	return nil, nil
}

func saveToESAndDB(reqId string, userInfo *login.RedisUserInfo, req *RentalInfoCreateRequest) error {
	var (
		obfsTitle string
		obfsDesc  string

		dbData db_models.RentalInfo
		esData es_model.RentalInfo

		wg = errgroup.Group{}
	)

	wg.Go(func() error {
		// 混淆数据
		for _, item := range req.Title {
			obfsTitle += "&#x" + strconv.FormatInt(int64(item), 16) + ";"
		}

		for _, item := range req.Desc {
			obfsDesc += "&#x" + strconv.FormatInt(int64(item), 16) + ";"
		}

		priceIndex := make([]float64, 0)
		temp := req.Price
		for temp != 0 {
			lastNum := temp % 10
			priceIndex = append(priceIndex, numMapMovement[lastNum])
			temp = temp / 10
		}

		reverse(priceIndex)

		esData = es_model.RentalInfo{
			RequestId:     reqId,
			Uid:           userInfo.Uid,
			Title:         req.Title,
			ObfsTitle:     obfsTitle,
			Cover:         req.Cover,
			Pics:          req.Pics,
			Area:          req.Area,
			Price:         priceIndex,
			RentAvailTime: req.RentAvailTime,
			RentTermFrom:  req.RentTermFrom,
			RentTermTo:    req.RentTermTo,
			Province:      req.Province,
			City:          req.City,
			Location:      req.Location,
			Desc:          req.Desc,
			ObfsDesc:      obfsDesc,
			Tags:          req.Tags,
			HouseType:     req.HouseType,
			RoomType:      req.RoomType,
			Furniture:     req.Furniture,
			Type:          req.Type,
			Rooms:         req.Rooms,
		}

		return nil
	})

	wg.Go(func() error {

		pics, _ := json.Marshal(req.Pics)
		tags, _ := json.Marshal(req.Tags)
		roomTypes, _ := json.Marshal(req.RoomType)
		furniture, _ := json.Marshal(req.Furniture)
		rooms, _ := json.Marshal(req.Rooms)

		dbData = db_models.RentalInfo{
			Uid:           userInfo.Uid,
			Title:         req.Title,
			Cover:         req.Cover,
			Pics:          string(pics),
			Area:          req.Area,
			Price:         req.Price,
			RentAvailTime: req.RentAvailTime,
			RentTermFrom:  req.RentTermFrom,
			RentTermTo:    req.RentTermTo,
			Province:      req.Province,
			City:          req.City,
			Location:      req.Location,
			Desc:          req.Desc,
			Tags:          string(tags),
			HouseType:     req.HouseType,
			RoomType:      string(roomTypes),
			Furniture:     string(furniture),
			Type:          db_models.RentTypeStrMapInt[req.Type],
			Rooms:         string(rooms),
		}

		return config.GetMysqlCli().Create(&dbData).Error
	})

	err := wg.Wait()
	if err != nil {
		return err
	}

	esData.ID = dbData.ID
	rentalInfoChannel <- esData

	return nil
}

func reverse(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}
