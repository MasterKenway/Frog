package log

import (
	"fmt"
	"os"

	"frog/module/common/tools"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logger *zap.SugaredLogger
)

func init() {
	cfg := zap.NewDevelopmentConfig()
	cfg.DisableStacktrace = true
	cfg.Encoding = "json"
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stderr"}

	rawLogger, err := cfg.Build(zap.AddCallerSkip(1), zap.WrapCore(zapCore))
	if err != nil {
		panic(err)
	}
	logger = rawLogger.Sugar()

	go func() {
		select {
		case <-tools.Done:
			logger.Sync()
		}
	}()
}

func zapCore(c zapcore.Core) zapcore.Core {
	curDir, _ := os.Getwd()
	logFilePath := curDir + string(os.PathSeparator) + "main.log"
	_, err := os.Stat(logFilePath)
	if err != nil {
		_, err = os.Create(logFilePath)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(logFilePath)

	//if runtime.GOOS == "windows" {
	//	logFilePath = "winfile:///" + logFilePath
	//}
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    500, // megabytes
		MaxBackups: 30,
		MaxAge:     30, // days
	})

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
		w,
		zap.DebugLevel,
	)
	cores := zapcore.NewTee(c, core)

	return cores
}
func Info(a ...interface{}) {
	logger.Info(a)
}

func Infof(format string, a ...interface{}) {
	logger.Infof(format, a...)
}

func Debug(a ...interface{}) {
	logger.Debug(a)
}

func Debugf(format string, a ...interface{}) {
	logger.Debugf(format, a...)
}

func Error(a ...interface{}) {
	logger.Error(a)
}

func Errorf(format string, a ...interface{}) {
	logger.Errorf(format, a...)
}

func Warn(a ...interface{}) {
	logger.Warn(a)
}

func Warnf(format string, a ...interface{}) {
	logger.Warnf(format, a...)
}
