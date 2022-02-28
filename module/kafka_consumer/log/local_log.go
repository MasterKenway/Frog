package log

import (
	"os"

	"frog/module/common"

	"go.uber.org/zap"
)

var (
	logger *zap.SugaredLogger
)

func init() {
	curDir, _ := os.Getwd()
	logFilePath := curDir + string(os.PathSeparator) + "main.log"
	_, err := os.Stat(logFilePath)
	if err != nil {
		_, err = os.Create(logFilePath)
		if err != nil {
			panic(err)
		}
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Encoding = "json"
	cfg.OutputPaths = []string{"stdout", logFilePath}
	cfg.ErrorOutputPaths = []string{"stderr", logFilePath}

	rawLogger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	logger = rawLogger.Sugar()

	go func() {
		select {
		case <-common.Done:
			logger.Sync()
		}
	}()
}

func Info(a ...interface{}) {
	logger.Info(a)
}

func Infof(format string, a ...interface{}) {
	logger.Infof(format, a)
}

func Debug(a ...interface{}) {
	logger.Debug(a)
}

func Debugf(format string, a ...interface{}) {
	logger.Debugf(format, a)
}

func Error(a ...interface{}) {
	logger.Error(a)
}

func Errorf(format string, a ...interface{}) {
	logger.Errorf(format, a)
}

func Warn(a ...interface{}) {
	logger.Warn(a)
}

func Warnf(format string, a ...interface{}) {
	logger.Warnf(format, a)
}
