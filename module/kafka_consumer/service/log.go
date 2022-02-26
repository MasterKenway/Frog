package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"sync"
	"time"

	"frog/module/common"
	"frog/module/common/constant"
	"frog/module/common/model/api_models"
	"frog/module/common/model/es_model"
	"frog/module/kafka_consumer/config"
	"frog/module/kafka_consumer/log"

	"github.com/Shopify/sarama"
)

var (
	consumeChan = make(chan struct{}, 10)
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
			msgs := messages
			messages = make([]*sarama.ConsumerMessage, 0)
			lock.Unlock()
			consume(msgs)
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
	var (
		ctx     = context.Background()
		rawLogs = make([]api_models.RawLog, 0)
	)

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

	indexer, err := config.GetESIndexer(es_model.ESLog{}.Index())
	if err != nil {
		log.Errorf("failed to get es indexer, %s", err.Error())
	}

	for _, esLog := range esLogs {
		jsonData, _ := json.Marshal(esLog)
		err := indexer.Add(ctx, esutil.BulkIndexerItem{
			Action: "index",
			Body:   bytes.NewReader(jsonData),
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, respItem esutil.BulkIndexerResponseItem, err error) {
				log.Errorf("bulk index failed, result: %s", respItem.Result)
				if err != nil {
					log.Errorf("bulk index failed, err: %s", err.Error())
				}
			},
		})
		if err != nil {
			log.Errorf("indexer.Add %s", err.Error())
			continue
		}
	}

	err = indexer.Close(ctx)
	if err != nil {
		log.Errorf("indexer.Close %s", err.Error())
		return
	}

	stats := indexer.Stats()
	if stats.NumFailed > 0 {
		log.Errorf("stats.NumFailed %d", stats.NumFailed)
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
