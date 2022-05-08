package main

import (
	"os"
	"os/signal"
	"syscall"

	"frog/module/common/tools"
	"frog/module/kafka_consumer/internel/log"
	"frog/module/kafka_consumer/internel/service"

	"github.com/jasonlvhit/gocron"
)

func main() {

	log.Info("Start server...")

	_ = tools.CronService.Every(1).Day().From(gocron.NextTick()).Do(service.DeleteLog)
	go func() {
		<-tools.CronService.Start()
	}()

	service.ConsumeLog()
	go service.ConsumeLoop()

	log.Info("Server Stated...")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")
	close(tools.Done)
}
