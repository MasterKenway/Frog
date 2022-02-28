package middleware

import (
	"bytes"
	"io/ioutil"
	"strings"

	"frog/module/common/constant"
	"frog/module/main_service/internal/tools"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID(ctx *gin.Context) {
	ctx.Set(constant.CtxKeyRequestID, uuid.New().String())

	parts := strings.Split(ctx.Request.RemoteAddr, ":")
	ctx.Set(constant.CtxKeyRemoteIP, parts[0])
	ctx.Set(constant.CtxKeyRemotePort, parts[1])

	reqBody, err := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body.Close()
	if err != nil {
		tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeInternalError, constant.MsgInternalError)
		return
	}

	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(reqBody))

	cmd, err := jsonparser.GetString(reqBody, "Cmd")
	if err != nil {
		tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeBadRequest, constant.MsgParamInvalid)
		return
	}

	ctx.Set(constant.CtxKeyCmd, cmd)
	ctx.Set(constant.CtxKeyReqBody, reqBody)

	ctx.Next()
}
