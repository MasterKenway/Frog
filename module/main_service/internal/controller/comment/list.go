package comment

import (
	"frog/module/common/constant"
	"frog/module/common/model/api_models"
	"frog/module/common/model/db_models"
	"frog/module/main_service/internal/config"
	"frog/module/main_service/internal/log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCommentListController() api_models.ApiInterface {
	return &CommentListRequest{}
}

type CommentListRequest struct {
	ID int `json:"ID"`
}

func (req *CommentListRequest) GetResult(ctx *gin.Context) (interface{}, *api_models.APIError) {
	var (
		reqId  = ctx.GetString(constant.CtxKeyRequestID)
		dbData = make([]db_models.Comment, 0)
	)

	err := config.GetMysqlCli().Model(&db_models.Comment{}).Where("rental_id = ?", req.ID).Scan(&dbData).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Errorf(reqId, err.Error())
		return nil, &api_models.APIError{
			Code:    constant.CodeInternalError,
			Message: constant.MsgInternalError,
		}
	}

	return dbData, nil
}
