package middleware

import (
	"context"
	"strings"

	"graduation-project/module/main_service/internal/config"
	"graduation-project/module/main_service/internal/constant"
	"graduation-project/module/main_service/internal/log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func AntiBlackIPs(ctx *gin.Context) {
	var (
		reqId = ctx.GetString(constant.CtxKeyRequestID)
	)

	parts := strings.Split(ctx.Request.RemoteAddr, ":")
	ctx.Set(constant.CtxKeyRemoteIP, parts[0])
	ctx.Set(constant.CtxKeyRemotePort, parts[1])

	stamp, err := config.GetRedisCli().Get(context.Background(), constant.RedisKeyIPStamp).Bytes()
	if err != nil && err != redis.Nil {
		log.Errorf(reqId, "failed get ip stamp")
	} else {
		if len(stamp) == 1 {
			if (stamp[0]>>(7-constant.IPStampProxy))&1 == 1 {
				ctx.Set(constant.CtxKeyIsProxy, true)
			}

			if (stamp[0]>>(7-constant.IPStampSpider))&1 == 1 || (stamp[0]>>(7-constant.IPStampBot))&1 == 1 {
				ctx.Set(constant.CtxKeyIsBot, true)
			}

			if (stamp[0]>>(7-constant.IPStampQuickConn))&1 == 1 {
				ctx.Set(constant.CtxKeyIsQuickConn, true)
			}
		}
	}

	ctx.Next()
}
