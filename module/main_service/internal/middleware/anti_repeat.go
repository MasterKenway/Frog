package middleware

import (
	"context"
	"strconv"
	"time"

	"graduation-project/module/main_service/internal/config"
	"graduation-project/module/main_service/internal/constant"
	"graduation-project/module/main_service/internal/log"
	"graduation-project/module/main_service/internal/tools"

	"github.com/gin-gonic/gin"
)

func AntiRepeat(ctx *gin.Context) {
	var (
		reqId = ctx.GetString(constant.CtxKeyRequestID)
	)

	if ts := ctx.GetHeader(constant.HeaderKeyTimeStamp); ts != "" {
		tsInt, err := strconv.ParseInt(ts, 10, 64)
		if err != nil {
			log.Errorf(reqId, "failed to parse timestamp, %s", err.Error())
			tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeBadRequest, constant.MsgTimeStampOutdated)
			return
		}
		if time.Unix(tsInt, 0).Sub(time.Now()) > time.Minute {
			log.Errorf(reqId, "timestamp outdated")
			tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeBadRequest, constant.MsgTimeStampOutdated)
			return
		}
	} else {
		log.Errorf(reqId, "timestamp invalid")
		tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeBadRequest, constant.MsgParamInvalid)
		return
	}

	if nonce := ctx.GetHeader(constant.HeaderKeyNonce); nonce != "" {
		err := config.GetRedisCli().Set(context.Background(), constant.RedisKeyNonce+nonce, true, 5*time.Minute).Err()
		if err != nil {
			log.Errorf(reqId, "redis set %s, err: %s", constant.RedisKeyNonce+nonce, err.Error())
		}
	} else {
		log.Errorf(ctx.GetString(constant.CtxKeyRequestID), "nonce invalid")
		tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeBadRequest, constant.MsgParamInvalid)
		return
	}

	ctx.Next()
}
