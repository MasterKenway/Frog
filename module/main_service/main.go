package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"frog/module/common/constant"
	"frog/module/common/tools"
	"frog/module/main_service/internal/controller"
	"frog/module/main_service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	go func() {
		<-tools.CronService.Start()
	}()

	r := gin.New()

	r.Use(
		middleware.Logger,
		middleware.RequestID,
		middleware.AntiRepeat,
		middleware.AntiBlackIPs,
		middleware.AntiUA,
		middleware.RateLimit,
		middleware.Captcha,
		middleware.ValidateLogin,
	)

	r.POST(constant.InterfaceEntry, controller.InterfaceHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	close(tools.Done)
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")

}
