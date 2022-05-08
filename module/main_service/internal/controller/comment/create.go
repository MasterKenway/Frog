package comment

import (
	"database/sql"
	"frog/module/common/constant"
	"frog/module/common/model/api_models"
	"frog/module/common/model/db_models"
	"frog/module/main_service/internal/config"
	"frog/module/main_service/internal/controller/auth/login"
	"frog/module/main_service/internal/log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCommentCreateController() api_models.ApiInterface {
	return &CommentCreateRequest{}
}

type CommentCreateRequest struct {
	RentalId int    `json:"RentalId" validate:"required"`
	Content  string `json:"Content" validate:"required"`
}

func (req *CommentCreateRequest) GetResult(ctx *gin.Context) (interface{}, *api_models.APIError) {
	var (
		reqId       = ctx.GetString(constant.CtxKeyRequestID)
		rentalInfo  db_models.RentalInfo
		userInfo, _ = ctx.Get(constant.CtxKeyUserInfo)

		notifs      = make([]db_models.Notification, 0)
		notifyUsers []string
	)

	userInfoObject := userInfo.(*login.RedisUserInfo)

	err := config.GetMysqlCli().Model(db_models.RentalInfo{}).Where("id = ?", req.RentalId).First(&rentalInfo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &api_models.APIError{
				Code:    constant.CodeBadRequest,
				Message: "Rental Info Not Found",
			}
		}
		log.Errorf(reqId, "%s", err.Error())
		return nil, &api_models.APIError{
			Code:    constant.CodeInternalError,
			Message: constant.MsgInternalError,
		}
	}

	dbModel := db_models.Comment{
		RentalId:   uint(req.RentalId),
		PublishUid: rentalInfo.Uid,
		Uid:        userInfoObject.Uid,
		Username:   userInfoObject.Username,
		Content:    req.Content,
	}

	err = config.GetMysqlCli().Create(&dbModel).Error
	if err != nil {
		log.Errorf(reqId, "%s", err.Error())
		return nil, &api_models.APIError{
			Code:    constant.CodeBadRequest,
			Message: "Save to db failed",
		}
	}

	err = config.GetMysqlCli().Model(&db_models.Comment{}).Select("DISTINCT uid").
		Where("rental_id = ?", rentalInfo.ID).
		Where("rental_id != ?", rentalInfo.Uid).
		Pluck("uid", &notifyUsers).Error
	if err != nil {
		log.Errorf(reqId, err.Error())
		return nil, &api_models.APIError{
			Code:    constant.CodeInternalError,
			Message: constant.MsgInternalError,
		}
	}

	notifyUsers = append(notifyUsers, rentalInfo.Uid)
	for _, uid := range notifyUsers {
		notifs = append(notifs, db_models.Notification{
			Type:      0,
			Uid:       uid,
			CommentId: sql.NullInt32{Int32: int32(dbModel.Id), Valid: true},
			RentalId:  sql.NullInt32{Int32: int32(dbModel.RentalId), Valid: true},
			Content:   req.Content,
			IsRead:    0,
		})
	}

	err = config.GetMysqlCli().Create(&notifs).Error
	if err != nil {
		log.Errorf(reqId, err.Error())
		return nil, &api_models.APIError{
			Code:    constant.CodeInternalError,
			Message: constant.MsgInternalError,
		}
	}

	return nil, nil
}
