package config

type KafkaConfig struct {
	Endpoint []string          `json:"endpoint"`
	GroupID  string            `json:"group_id"`
	Topics   map[string]string `json:"topics"`
}
