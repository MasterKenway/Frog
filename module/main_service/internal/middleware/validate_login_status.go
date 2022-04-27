package middleware

import (
	"context"
	"encoding/json"
	"time"

	"frog/module/common/constant"
	"frog/module/main_service/internal/config"
	"frog/module/main_service/internal/controller/auth/login"
	"frog/module/main_service/internal/log"
	"frog/module/main_service/internal/tools"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var (
	LoginRequiredUrl = map[string]bool{
		"CreateRentalInfo": true,
	}
)

func ValidateLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, ok := LoginRequiredUrl[ctx.GetString(constant.CtxKeyCmd)]; !ok {
			ctx.Next()
			return
		}

		cookie, err := ctx.Cookie(constant.CookieKeyLoginCert)
		if err != nil {
			tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeLoginRequired, constant.MsgNotLogin)
			return
		}

		result, err := config.GetRedisCli().Get(context.Background(), tools.GetRedisKeyLoginCert(cookie)).Bytes()
		if err != nil {
			if err == redis.Nil {
				tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeLoginRequired, constant.MsgNotLogin)
			} else {
				log.Errorf(ctx.GetString(constant.CtxKeyRequestID), "failed to get user info from redis, %s", err.Error())
				tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeInternalError, constant.MsgInternalError)
			}
			return
		}

		userInfo := &login.RedisUserInfo{}
		err = json.Unmarshal(result, userInfo)
		if err != nil {
			log.Errorf(ctx.GetString(constant.CtxKeyRequestID), "failed to unmarshal user info, %s", err.Error())
			tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeLoginRequired, constant.MsgNotLogin)
			return
		}

		// 登录态有效则每次请求延长登录有效时间
		config.GetRedisCli().Expire(context.Background(), tools.GetRedisKeyLoginCert(cookie), 30*time.Minute)

		ctx.Set(constant.CtxKeyUserInfo, userInfo)
		ctx.Next()
	}
}
