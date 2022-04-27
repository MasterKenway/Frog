package middleware

import (
	"context"
	"regexp"
	"time"

	"frog/module/common/constant"
	"frog/module/common/model/db_models"
	comTools "frog/module/common/tools"
	"frog/module/main_service/internal/config"
	"frog/module/main_service/internal/log"
	"frog/module/main_service/internal/tools"

	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
)

const (
	uaDefaultConfig = "PostmanRuntime|headless"
)

var (
	uaReg   *regexp.Regexp
	uaCache map[string]bool
)

func init() {
	err := comTools.CronService.Every(1).Hour().From(gocron.NextTick()).Do(syncUAConfig)
	if err != nil {
		log.Errorf("gocron-do", "failed to add cron job, %s", err.Error())
	}

	err = comTools.CronService.Every(1).Hour().From(gocron.NextTick()).Do(func() { uaCache = map[string]bool{} })
	if err != nil {
		log.Errorf("gocron-do", "failed to add cron job, %s", err.Error())
	}
}

func AntiUA() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqId := ctx.GetString(constant.CtxKeyRequestID)
		ua := ctx.GetHeader("user-agent")

		if ua == "" {
			tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeForbidden, "UA Invalid")
			return
		}

		if _, ok := uaCache[ua]; ok {
			log.Warnf(reqId, "Hit UA black list")
			config.GetRedisCli().Set(context.Background(), tools.GetRedisKeyFingerPrint(ctx.GetString(constant.RedisKeyBlockFP)), 1, 15*time.Minute)
			tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeForbidden, constant.MsgIllegalRequest)
			return
		}

		if uaReg.Match([]byte(ua)) {
			log.Warnf(reqId, "Hit UA black list")
			uaCache[ua] = true
			config.GetRedisCli().Set(context.Background(), tools.GetRedisKeyFingerPrint(ctx.GetString(constant.RedisKeyBlockFP)), 1, 15*time.Minute)
			tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeForbidden, constant.MsgIllegalRequest)
			return
		}

		ctx.Next()
	}
}

func syncUAConfig() {
	var uaConfig db_models.UaConfig
	err := config.GetMysqlCli().Model(db_models.UaConfig{}).Where("is_delete = 0").First(&uaConfig).Error
	if err != nil {
		log.Errorf("sync-ua-config", "failed to query config from db, %s", err.Error())
		uaReg, _ = regexp.Compile(uaDefaultConfig)
		return
	}

	uaRegTemp, err := regexp.Compile(uaConfig.Ua)
	if err != nil {
		log.Errorf("sync-ua-config", "failed to compile regex, %s", err.Error())
		uaReg, _ = regexp.Compile(uaDefaultConfig)
		return
	}

	uaReg = uaRegTemp
}
