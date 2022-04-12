package service

import (
	"bytes"
	"context"
	"encoding/json"
	"frog/module/common/model/db_models"
	"frog/module/common/tools"
	config2 "frog/module/kafka_consumer/internel/config"
	"frog/module/kafka_consumer/internel/log"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"sync"
	"time"

	"frog/module/common/constant"
	"frog/module/common/model/api_models"
	"frog/module/common/model/es_model"
	"github.com/Shopify/sarama"
)

var (
	consumeChan = make(chan struct{}, 10)
	messages    = make([][]byte, 0)
	lock        = sync.Mutex{}
)

func ConsumeLog() {
	log.Info("creating consume group successfully")
	kafkaConsumer := *config2.GetKafkaConsumer()
	claimConsumer := LogConsumer{}
	go func() {
		for {
			err := kafkaConsumer.Consume(context.Background(), []string{config2.GetKafkaTopic(constant.KafkaKeyLogTopic)}, &claimConsumer)
			if err != nil {
				log.Errorf("failed to consume kafka topic, %s", err.Error())
				continue
			}
		}

	}()
	log.Info("create consume group successfully")
}

func ConsumeLoop() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-tools.Done:
			ticker.Stop()
			lock.Lock()
			msgs := messages
			messages = make([][]byte, 0)
			lock.Unlock()
			consume(msgs)
			return
		case <-ticker.C:
			lock.Lock()
			msgs := messages
			messages = make([][]byte, 0)
			lock.Unlock()
			consume(msgs)
		case <-consumeChan:
			lock.Lock()
			msgs := messages
			messages = make([][]byte, 0)
			lock.Unlock()
			consume(msgs)
		}
	}
}

func consume(msgs [][]byte) {
	var (
		ctx     = context.Background()
		rawLogs = make([]api_models.RawLog, 0)
	)

	if len(msgs) == 0 {
		return
	}

	for _, msg := range msgs {
		var rawLog api_models.RawLog
		err := json.Unmarshal(msg, &rawLog)
		if err != nil {
			log.Errorf("failed to unmarshal log message, %s, rawLog: %s", err.Error(), string(msg))
			continue
		}
		rawLogs = append(rawLogs, rawLog)
	}

	esLogs := make([]es_model.ESLog, 0)
	dbLogs := make([]db_models.Log, 0)
	for _, rawLog := range rawLogs {
		esLogs = append(esLogs, es_model.ESLog{
			Time:      rawLog.Time,
			Level:     rawLog.Level,
			Caller:    rawLog.Caller,
			RequestID: rawLog.RequestID,
			Message:   rawLog.Message,
		})

		dbLogs = append(dbLogs, db_models.Log{
			Time:      rawLog.Time,
			Level:     rawLog.Level,
			Caller:    rawLog.Caller,
			RequestId: rawLog.RequestID,
			Message:   rawLog.Message,
		})
	}

	err := config2.GetMysqlCli().Create(&dbLogs).Error
	if err != nil {
		log.Errorf("failed to save log to db, %s", err.Error())
	}

	indexer, err := config2.GetESIndexer(es_model.ESLog{}.Index())
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
	for message := range claim.Messages() {
		lock.Lock()
		messages = append(messages, message.Value)
		if len(messages) > 100 {
			consumeChan <- struct{}{}
		}
		lock.Unlock()
		session.MarkMessage(message, "")
	}

	return nil
}
