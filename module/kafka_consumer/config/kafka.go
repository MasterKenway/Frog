package config

import (
	"encoding/json"

	"graduation-project/module/common/config"
	"graduation-project/module/common/constant"

	"github.com/Shopify/sarama"
)

var (
	KafkaConsumer *sarama.ConsumerGroup
	KafkaConf     *config.KafkaConfig
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

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := sarama.NewConsumerGroup(KafkaConf.Endpoint, KafkaConf.GroupID, kafkaConfig)
	KafkaConsumer = &consumer

	return KafkaConsumer
}

func GetKafkaTopic(topicKey string) string {
	return KafkaConf.Topics[topicKey]
}
