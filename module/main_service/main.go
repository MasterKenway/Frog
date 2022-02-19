package main

import (
	"graduation-project/module/main_service/internal/constant"
	"graduation-project/module/main_service/internal/controller"
	"graduation-project/module/main_service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(middleware.RequestID)
	r.Use(middleware.AntiRepeat)
	r.Use(middleware.AntiBlackIPs)
	r.Use(middleware.AntiUA)
	r.Use(middleware.GateWay)
	r.Use(middleware.Captcha)

	r.POST(constant.InterfaceEntry, controller.InterfaceHandler)
}
