package middleware

import (
	"context"

	"frog/module/common/constant"
	"frog/module/main_service/internal/config"
	"frog/module/main_service/internal/tools"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func BlockFingerPrint() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, err := config.GetRedisCli().Exists(context.Background(), tools.GetRedisKeyFingerPrint(ctx.GetString(constant.CtxKeyFingerPrint))).Uint64()
		if err != nil {
			if err != redis.Nil {
				tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeInternalError, constant.MsgInternalError)
				return
			} else {
				ctx.Next()
			}
			return
		}

		if result == 0 {
			tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeForbidden, constant.MsgIllegalRequest)
			return
		}

		ctx.Next()
		return
	}
}
