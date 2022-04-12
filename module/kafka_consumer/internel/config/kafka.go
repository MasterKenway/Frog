package config

import (
	"encoding/json"
	"frog/module/kafka_consumer/internel/log"
	log2 "log"
	"os"

	"frog/module/common/config"
	"frog/module/common/constant"

	"github.com/Shopify/sarama"
)

var (
	KafkaConsumer *sarama.ConsumerGroup
	KafkaConf     = &config.KafkaConfig{}
)

func GetKafkaConsumer() *sarama.ConsumerGroup {
	if KafkaConsumer != nil {
		return KafkaConsumer
	}

	kafkaEtcdConfBytes, err := GetConfig(constant.EtcdKeyKafkaConfig)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(kafkaEtcdConfBytes, KafkaConf)
	if err != nil {
		panic(err)
	}

	log.Infof("%+v", *KafkaConf)

	kafkaConfig := sarama.NewConfig()
	sarama.Logger = log2.New(os.Stdout, "[sarama] ", log2.LstdFlags)
	kafkaConfig.Net.SASL.Enable = false
	version, err := sarama.ParseKafkaVersion("3.1.0")
	if err != nil {
		panic(err)
	}
	kafkaConfig.Version = version
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := sarama.NewConsumerGroup(KafkaConf.Endpoint, KafkaConf.GroupID, kafkaConfig)
	if err != nil {
		panic(err)
	}
	KafkaConsumer = &consumer

	return KafkaConsumer
}

func GetKafkaTopic(topicKey string) string {
	return KafkaConf.Topics[topicKey]
}
