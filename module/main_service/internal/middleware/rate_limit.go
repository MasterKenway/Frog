package middleware

import (
	"context"
	"time"

	"frog/module/common/constant"
	"frog/module/main_service/internal/config"
	"frog/module/main_service/internal/log"
	comTool "frog/module/main_service/internal/tools"

	"github.com/gin-gonic/gin"
)

func RateLimit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		redisKey := comTool.GetRedisKeyRateLimit(ctx.GetString(constant.CtxKeyRemoteIP))
		rate, err := config.GetRedisCli().Incr(context.Background(), redisKey).Uint64()
		if err != nil {
			log.Errorf(ctx.GetString(constant.CtxKeyRequestID), "failed to get time limit, %s", err.Error())
			ctx.Next()
			return
		}

		if rate == 1 {
			config.GetRedisCli().Expire(context.Background(), redisKey, 5*time.Minute)
		}

		if rate > config.GetRateLimitConfig().TimesPerSec {
			comTool.CtxAbortWithCodeAndMsg(ctx, constant.CodeRateLimit, constant.MsgRateLimit)
			return
		}

		ctx.Next()
	}
}
