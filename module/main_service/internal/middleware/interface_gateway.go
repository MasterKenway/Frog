package middleware

import (
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"graduation-project/module/main_service/internal/constant"
	"graduation-project/module/main_service/internal/tools"

	"github.com/buger/jsonparser"
)

var (
	LoginNeededUrl = map[string]bool{}
)

func GateWay(ctx *gin.Context) {
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
