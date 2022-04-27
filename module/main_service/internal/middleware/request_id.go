package middleware

import (
	"bytes"
	"frog/module/main_service/internal/log"
	"io/ioutil"

	"frog/module/common/constant"
	comTools "frog/module/common/tools"
	"frog/module/main_service/internal/tools"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
)

func RequestID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(constant.CtxKeyRemoteIP, comTools.GetRemoteAddr(ctx))

		reqBody, err := ioutil.ReadAll(ctx.Request.Body)
		ctx.Request.Body.Close()
		if err != nil {
			tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeInternalError, constant.MsgInternalError)
			return
		}

		ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(reqBody))

		if ctx.Request.RequestURI == constant.InterfaceEntry {
			cmd, err := jsonparser.GetString(reqBody, "Cmd")
			if err != nil {
				log.Errorf(ctx.GetString(constant.CtxKeyRequestID), "Cmd Not Found")
				tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeBadRequest, constant.MsgParamInvalid)
				return
			}

			ctx.Set(constant.CtxKeyCmd, cmd)
		}

		fp := ctx.Request.Header.Get(constant.HeaderKeyFingerPrint)

		ctx.Set(constant.CtxKeyReqBody, reqBody)
		ctx.Set(constant.CtxKeyFingerPrint, fp)

		ctx.Next()
	}
}
