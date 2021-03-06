package log

import (
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
	//Producer sarama.AsyncProducer
	Producer sarama.SyncProducer
	Topic    string
}

func (k KafkaWriter) Write(p []byte) (n int, err error) {
	err = k.Producer.SendMessages([]*sarama.ProducerMessage{{Topic: k.Topic, Key: nil, Value: sarama.ByteEncoder(p), Partition: 0}})
	//k.Producer.Input() <- &sarama.ProducerMessage{Topic: k.Topic, Key: nil, Value: sarama.ByteEncoder(p), Partition: 0}
	return len(p), err
}

type customLogger struct {
	log *zap.SugaredLogger
}

func init() {
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.FatalLevel
	})

	kw := KafkaWriter{
		//Producer: config.GetKafkaAsyncProducer(),
		Producer: config.GetKafkaSyncProducer(),
		Topic:    config.GetKafkaTopic(constant.KafkaKeyLogTopic),
	}
	topicError := zapcore.AddSync(kw)
	encodeConfig := zap.NewDevelopmentEncoderConfig()
	encodeConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	kafkaEncoder := zapcore.NewJSONEncoder(encodeConfig)
	kafkaCore := zapcore.NewCore(kafkaEncoder, topicError, lowPriority)
	core := zapcore.NewTee(kafkaCore)
	logger = customLogger{log: zap.New(core).WithOptions(zap.AddCallerSkip(1)).WithOptions(zap.AddCaller()).Sugar()}
}

func Info(reqId string, a ...interface{}) {
	logger.log.With("RID", reqId).Info(a...)
}

func Infof(reqId string, format string, a ...interface{}) {
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

func Fatal(reqId string, a ...interface{}) {
	logger.log.With("RID", reqId).Fatal(a...)
}

func Fatalf(reqId, format string, a ...interface{}) {
	logger.log.With("RID", reqId).Fatalf(format, a...)
}
