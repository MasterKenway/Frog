package middleware

import (
	"context"
	"encoding/json"

	"frog/module/common/constant"
	"frog/module/main_service/internal/config"
	"frog/module/main_service/internal/controller/auth/login"
	"frog/module/main_service/internal/log"
	"frog/module/main_service/internal/tools"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var (
	LoginRequiredUrl = map[string]bool{}
)

func ValidateLogin(ctx *gin.Context) {
	if _, ok := LoginRequiredUrl[ctx.GetString(constant.CtxKeyCmd)]; !ok {
		ctx.Next()
		return
	}

	cookie, err := ctx.Cookie(constant.CookieKeyLoginCert)
	if err != nil {
		tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeLoginRequired, constant.MsgNotLogin)
		return
	}

	result, err := config.GetRedisCli().Get(context.Background(), cookie).Bytes()
	if err != nil {
		if err == redis.Nil {
			tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeLoginRequired, constant.MsgNotLogin)
		} else {
			log.Errorf(ctx.GetString(constant.CtxKeyRequestID), "failed to get user info from redis, %s", err.Error())
			tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeInternalError, constant.MsgInternalError)
		}
		return
	}

	var userInfo login.RedisUserInfo
	err = json.Unmarshal(result, &userInfo)
	if err != nil {
		log.Errorf(ctx.GetString(constant.CtxKeyRequestID), "failed to unmarshal user info, %s", err.Error())
		tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeLoginRequired, constant.MsgNotLogin)
		return
	}

	ctx.Set(constant.CtxKeyUserInfo, userInfo)
	ctx.Next()
}
