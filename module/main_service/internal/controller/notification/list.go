package notification

import (
	"frog/module/common/constant"
	"frog/module/common/model/api_models"
	"frog/module/common/model/db_models"
	"frog/module/main_service/internal/config"
	"frog/module/main_service/internal/controller/auth/login"
	"frog/module/main_service/internal/log"

	"github.com/gin-gonic/gin"
)

func GetNotificationListController() api_models.ApiInterface {
	return &NotificationListRequest{}
}

type NotificationListRequest struct {
}

func (req *NotificationListRequest) GetResult(ctx *gin.Context) (interface{}, *api_models.APIError) {
	var (
		reqId = ctx.GetString(constant.CtxKeyRequestID)

		userInfo, _    = ctx.Get(constant.CtxKeyUserInfo)
		userInfoObject = userInfo.(*login.RedisUserInfo)

		res []db_models.Notification
	)

	err := config.GetMysqlCli().Model(db_models.Notification{}).
		Where("uid = ?", userInfoObject.Uid).
		Where("is_read = 0").
		Scan(&res).Error
	if err != nil {
		log.Errorf(reqId, err.Error())
		return nil, &api_models.APIError{
			Code:    constant.CodeBadRequest,
			Message: err.Error(),
		}
	}

	return res, nil
}
