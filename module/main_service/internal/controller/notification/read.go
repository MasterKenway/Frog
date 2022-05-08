package notification

import (
	"frog/module/common/constant"
	"frog/module/common/model/api_models"
	"frog/module/common/model/db_models"
	"frog/module/main_service/internal/config"
	"frog/module/main_service/internal/log"
	"github.com/gin-gonic/gin"
)

func GetNotificationReadController() api_models.ApiInterface {
	return &NotificationReadRequest{}
}

type NotificationReadRequest struct {
	ID []int `json:"ID"`
}

func (req *NotificationReadRequest) GetResult(ctx *gin.Context) (interface{}, *api_models.APIError) {
	var (
		reqId = ctx.GetString(constant.CtxKeyRequestID)
	)

	err := config.GetMysqlCli().Model(&db_models.Notification{}).Where("id IN (?)", req.ID).Update("is_read", 1).Error
	if err != nil {
		log.Errorf(reqId, err.Error())
		return nil, &api_models.APIError{
			Code:    constant.CodeBadRequest,
			Message: err.Error(),
		}
	}

	return nil, nil
}
