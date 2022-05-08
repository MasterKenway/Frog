package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"frog/module/common/constant"
	"frog/module/common/model/api_models"
	"frog/module/common/tools"
	"frog/module/main_service/internal/controller"
	"frog/module/main_service/internal/log"
	"frog/module/main_service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	go func() {
		<-tools.CronService.Start()
	}()

	r := gin.New()

	r.Use(
		gin.CustomRecovery(CustomRecovery()),
		middleware.Logger(),
		middleware.RequestID(),
		middleware.AntiRepeat(),
		middleware.AntiBlackIPs(),
		middleware.AntiUA(),
		//middleware.RateLimit(),
		middleware.Captcha(),
		middleware.ValidateLogin(),
	)

	r.POST(constant.InterfaceUpload, controller.UploadToCos())
	r.POST(constant.InterfaceEntry, controller.InterfaceHandler)

	srv := &http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err.Error())
		}
	}()

	fmt.Println("Server Stated...")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")
	close(tools.Done)
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	fmt.Println("Server exiting")
}

func CustomRecovery() gin.RecoveryFunc {
	return func(ctx *gin.Context, err interface{}) {
		reqId := ctx.GetString(constant.CtxKeyRequestID)
		log.Errorf(reqId, "%v\n%s", err, string(debug.Stack()))
		ctx.JSON(http.StatusOK, api_models.APIResponse{ResponseInfo: api_models.ResponseInfo{
			Code:      constant.CodeInternalError,
			Error:     err,
			RequestID: reqId,
		}})
		ctx.Abort()
	}
}
