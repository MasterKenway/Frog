package controller

import (
	"encoding/json"
	"frog/module/common/constant"
	"frog/module/common/model/api_models"
	"frog/module/main_service/internal/log"
	"net/http"

	"frog/module/main_service/internal/tools"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func InterfaceHandler(ctx *gin.Context) {
	var (
		reqId = ctx.GetString(constant.CtxKeyRequestID)
	)

	defer ctx.Request.Body.Close()

	interf := ctx.GetString(constant.CtxKeyCmd)
	adapter, ok := ApiAdapter[interf]
	if !ok {
		tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeBadRequest, constant.MsgApiNotExists)
		return
	}

	reqBody, _ := ctx.Get(constant.CtxKeyReqBody)
	controllerBytes, _, _, err := jsonparser.Get(reqBody.([]byte), "Data")
	if err != nil {
		tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeBadRequest, constant.MsgParamInvalid)
		return
	}
	controller := adapter()
	err = json.Unmarshal(controllerBytes, &controller)
	if err != nil {
		tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeBadRequest, constant.MsgParamInvalid)
		return
	}

	err = Validate(controller)
	if err != nil {
		msg, _ := json.Marshal(api_models.APIResponse{ResponseInfo: api_models.ResponseInfo{
			RequestID: reqId,
			Code:      constant.CodeBadRequest,
			Message:   err.Error(),
		}})
		log.Infof(reqId, string(msg))
		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Writer.Write(msg)
		return
	}

	result, apiErr := controller.GetResult(ctx)
	if apiErr != nil {
		msg, _ := json.Marshal(api_models.APIResponse{ResponseInfo: api_models.ResponseInfo{
			RequestID: reqId,
			Code:      apiErr.Code,
			Message:   apiErr.Message,
			Error:     result,
		}})
		log.Infof(reqId, string(msg))
		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Writer.Write(msg)
	} else {
		msg, _ := json.Marshal(api_models.APIResponse{ResponseInfo: api_models.ResponseInfo{
			RequestID: reqId,
			Code:      constant.CodeSuccess,
			Message:   constant.MsgSuccess,
			Data:      result,
		}})
		log.Infof(reqId, string(msg))
		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Writer.Write(msg)
	}
}

var (
	valid = validator.New()
)

func Validate(apiInterface api_models.ApiInterface) error {
	return valid.Struct(apiInterface)
}
