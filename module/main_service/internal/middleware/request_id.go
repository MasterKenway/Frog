package middleware

import (
	"bytes"
	"io/ioutil"

	"frog/module/common/constant"
	comTools "frog/module/common/tools"
	"frog/module/main_service/internal/tools"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
)

func RequestID(ctx *gin.Context) {
	ctx.Set(constant.CtxKeyRemoteIP, comTools.GetRemoteAddr(ctx))

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
