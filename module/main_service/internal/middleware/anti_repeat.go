package middleware

import (
	"context"
	"frog/module/common/constant"
	"frog/module/common/model/db_models"
	"strconv"
	"time"

	"frog/module/main_service/internal/config"
	"frog/module/main_service/internal/log"
	"frog/module/main_service/internal/tools"

	"github.com/gin-gonic/gin"
)

func init() {
	GetEscapeIP()
}

var (
	escapeIPs = map[string]bool{}
)

func GetEscapeIP() {
	var escapes []db_models.EscapeMiddlewareIps
	err := config.GetMysqlCli().Model(db_models.EscapeMiddlewareIps{}).Where("is_delete = 0").Scan(&escapes).Error
	if err != nil {
		log.Errorf("get-escape-ips", "failed to query escape ips from db, %s", err.Error())
		return
	}

	for _, escape := range escapes {
		escapeIPs[escape.IP] = true
	}
}

func AntiRepeat(ctx *gin.Context) {
	var (
		reqId    = ctx.GetString(constant.CtxKeyRequestID)
		remoteIP = ctx.GetString(constant.CtxKeyRemoteIP)
	)

	if _, ok := escapeIPs[remoteIP]; ok {
		ctx.Next()
		return
	}

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
