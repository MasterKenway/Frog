package config

import "github.com/Shopify/sarama"

var (
	kafkaProducer *sarama.AsyncProducer
)

func GetKafkaProducer() sarama.AsyncProducer {
	if kafkaProducer != nil {
		return *kafkaProducer
	}

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Partitioner = sarama.NewRandomPartitioner
	kafkaConfig.Producer.Return.Errors = true
	p, err := sarama.NewAsyncProducer([]string{}, kafkaConfig)
	if err != nil {
		panic(err)
	}
	kafkaProducer = &p
	return *kafkaProducer
}
