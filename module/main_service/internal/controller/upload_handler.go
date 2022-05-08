package controller

import (
	"fmt"
	"net/http"

	"frog/module/common/constant"
	"frog/module/common/model/api_models"
	"frog/module/main_service/internal/config"
	"frog/module/main_service/internal/log"
	"frog/module/main_service/internal/tools"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tencentyun/cos-go-sdk-v5"
)

func UploadToCos() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			reqId = ctx.GetString(constant.CtxKeyRequestID)
		)

		file, err := ctx.FormFile("file")
		if err != nil {
			log.Errorf(reqId, err.Error())
			tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeBadRequest, constant.MsgParamInvalid)
			return
		}

		fileHeader, err := file.Open()
		if err != nil {
			log.Errorf(reqId, err.Error())
			tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeBadRequest, err.Error())
			return
		}

		filename := uuid.New().String()

		_, err = config.GetCOSClient().Object.Put(ctx, "upload/"+filename, fileHeader, &cos.ObjectPutOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{ContentType: file.Header.Get("Content-Type")},
		})
		if err != nil {
			log.Errorf(reqId, err.Error())
			tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeInternalError, err.Error())
			return
		}

		cosConfig := config.GetCOSConfig()
		ctx.JSON(http.StatusOK, api_models.APIResponse{ResponseInfo: api_models.ResponseInfo{
			Code:    constant.CodeSuccess,
			Message: constant.MsgSuccess,
			Data: map[string]string{
				"URL": fmt.Sprintf("https://%s.cos.ap-guangzhou.myqcloud.com/upload/%s", cosConfig.Bucket, filename),
			},
			RequestID: reqId,
		}})
	}
}
