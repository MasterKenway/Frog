package middleware

import (
	"context"
	"strconv"
	"time"

	"frog/module/common/constant"
	"frog/module/common/model/db_models"
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

func AntiRepeat() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
			result, err := config.GetRedisCli().GetSet(context.Background(), tools.GetRedisKeyNonce(nonce), 1).Int64()
			if err != nil {
				return
			}

			if result == 0 {
				config.GetRedisCli().Expire(context.Background(), tools.GetRedisKeyNonce(nonce), 5*time.Minute)
			} else {
				tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeForbidden, constant.MsgNonceInvalid)
				return
			}
		} else {
			log.Errorf(ctx.GetString(constant.CtxKeyRequestID), constant.MsgNonceInvalid)
			tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeBadRequest, constant.MsgParamInvalid)
			return
		}

		ctx.Next()
	}
}
