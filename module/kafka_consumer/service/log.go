package service

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"graduation-project/module/common"
	"graduation-project/module/common/constant"
	"graduation-project/module/common/model/api_models"
	"graduation-project/module/common/model/es_model"
	"graduation-project/module/kafka_consumer/config"
	"graduation-project/module/kafka_consumer/log"

	"github.com/Shopify/sarama"
	"github.com/olivere/elastic/v7"
)

var (
	consumeChan = make(chan struct{}, 3)
	messages    = make([]*sarama.ConsumerMessage, 0)
	lock        = sync.Mutex{}
)

func init() {
	ConsumeLog()
	go consumeLoop()
}

func ConsumeLog() {
	kafkaConsumer := *config.GetKafkaConsumer()
	claimConsumer := LogConsumer{}
	err := kafkaConsumer.Consume(context.Background(), []string{config.GetKafkaTopic(constant.KafkaKeyLogTopic)}, &claimConsumer)
	if err != nil {
		log.Errorf("failed to consume kafka topic, %s", err.Error())
		return
	}
}

func consumeLoop() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-common.Done:
			ticker.Stop()
			lock.Lock()
			msgs := messages
			messages = make([]*sarama.ConsumerMessage, 0)
			lock.Unlock()
			consume(msgs)
			return
		case <-ticker.C:
			lock.Lock()
			consumeChan <- struct{}{}
		case <-consumeChan:
			lock.Lock()
			msgs := messages
			messages = make([]*sarama.ConsumerMessage, 0)
			lock.Unlock()
			consume(msgs)
		}
	}
}

func consume(msgs []*sarama.ConsumerMessage) {
	rawLogs := make([]api_models.RawLog, 0)
	for _, msg := range msgs {
		var rawLog api_models.RawLog
		err := json.Unmarshal(msg.Value, &rawLog)
		if err != nil {
			log.Errorf("failed to unmarshal log message, %s, rawLog: %s", err.Error(), string(msg.Value))
			continue
		}
		rawLogs = append(rawLogs, rawLog)
	}

	esLogs := make([]es_model.ESLog, 0)
	for _, rawLog := range rawLogs {
		esLogs = append(esLogs, es_model.ESLog{
			Time:      rawLog.Time,
			Level:     rawLog.Level,
			Caller:    rawLog.Caller,
			RequestID: rawLog.RequestID,
			Message:   rawLog.Message,
		})
	}

	req := config.GetESCli().Bulk().Index(config.GetESIndexName(es_model.ESLog{}.Index()))
	for _, esLog := range esLogs {
		jsonData, _ := json.Marshal(esLog)
		doc := elastic.NewBulkIndexRequest().Doc(jsonData)
		req.Add(doc)
	}

	resp, err := req.Do(context.Background())
	if err != nil {
		log.Errorf("failed to save log to es, %s", err.Error())
		return
	}

	for _, item := range resp.Failed() {
		log.Errorf("resp.Failed: %s", item.Result)
	}
}

type LogConsumer struct {
}

func (l *LogConsumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (l *LogConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (l *LogConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := make([]*sarama.ConsumerMessage, 0)
	for message := range claim.Messages() {
		msgs = append(msgs, message)
	}

	lock.Lock()
	messages = append(messages, msgs...)
	if len(messages) > 100 {
		consumeChan <- struct{}{}
	}
	lock.Unlock()

	return nil
}
