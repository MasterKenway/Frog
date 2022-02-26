package config

import (
	"encoding/json"

	"frog/module/common/config"
	"frog/module/common/constant"

	"github.com/Shopify/sarama"
)

var (
	kafkaProducer *sarama.AsyncProducer
	KafkaConf     *config.KafkaConfig
)

func GetKafkaProducer() sarama.AsyncProducer {
	if kafkaProducer != nil {
		return *kafkaProducer
	}

	kafkaEtcdConfBytes, err := GetConfig(constant.EtcdKeyKafkaConfig)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(kafkaEtcdConfBytes, KafkaConf)
	if err != nil {
		panic(err)
	}

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	kafkaConfig.Producer.Return.Errors = true
	p, err := sarama.NewAsyncProducer(KafkaConf.Endpoint, kafkaConfig)
	if err != nil {
		panic(err)
	}
	kafkaProducer = &p
	return *kafkaProducer
}

func GetKafkaTopic(topicKey string) string {
	return KafkaConf.Topics[topicKey]
}
