package log

import (
	"fmt"
	"frog/module/common/constant"
	"frog/module/main_service/internal/config"

	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger customLogger
)

type KafkaWriter struct {
	Producer sarama.AsyncProducer
	Topic    string
}

func (k KafkaWriter) Write(p []byte) (n int, err error) {
	k.Producer.Input() <- &sarama.ProducerMessage{Topic: k.Topic, Key: nil, Value: sarama.ByteEncoder(p)}
	return len(p), nil
}

func (k *KafkaWriter) HandleErrors() {
	for {
		select {
		case err := <-k.Producer.Errors():
			fmt.Printf("k.Producer.Errors() %s\n", err.Error())
		}
	}
}

type customLogger struct {
	log *zap.SugaredLogger
}

func init() {
	// First, define our level-handling logic.
	//highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
	//	return lvl >= zapcore.ErrorLevel
	//})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	kw := KafkaWriter{
		Producer: config.GetKafkaProducer(),
		Topic:    config.GetKafkaTopic(constant.KafkaKeyLogTopic),
	}
	topicError := zapcore.AddSync(kw)
	kafkaEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
	kafkaCore := zapcore.NewCore(kafkaEncoder, topicError, lowPriority)
	core := zapcore.NewTee(kafkaCore)

	logger = customLogger{log: zap.New(core).WithOptions().Sugar()}
}

func Info(reqId string, a ...interface{}) {
	logger.log.With("RID", reqId).Info(a...)

}

func Infof(format, reqId string, a ...interface{}) {
	logger.log.With("RID", reqId).Infof(format, a...)
}

func Debug(reqId string, a ...interface{}) {
	logger.log.With("RID", reqId).Debug(a...)
}

func Debugf(reqId, format string, a ...interface{}) {
	logger.log.With("RID", reqId).Debugf(format, a...)
}

func Error(reqId string, a ...interface{}) {
	logger.log.With("RID", reqId).Error(a...)

}

func Errorf(reqId, format string, a ...interface{}) {
	logger.log.With("RID", reqId).Errorf(format, a...)
}

func Warn(reqId string, a ...interface{}) {
	logger.log.With("RID", reqId).Warn(a...)
}

func Warnf(reqId, format string, a ...interface{}) {
	logger.log.With("RID", reqId).Warnf(format, a...)
}
