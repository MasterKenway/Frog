package tools

import (
	"os"

	"frog/module/common/constant"

	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
)

var (
	Done = make(chan struct{})

	EnvType string

	CronService = gocron.NewScheduler()
)

func GetEnvType() string {
	envType, ok := os.LookupEnv("EnvType")
	if !ok || envType == "" {
		EnvType = "Testing"
	} else {
		EnvType = envType
	}

	return EnvType
}

func GetRemoteAddr(ctx *gin.Context) string {
	var ip string
	ip = ctx.GetHeader(constant.HeaderKeyXForwardedFor)
	if ip != "" {
		return ip
	}
	return ctx.ClientIP()
}
