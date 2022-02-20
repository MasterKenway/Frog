package config

type ElasticConfig struct {
	Urls     []string          `json:"url,omitempty"`
	Username string            `json:"username,omitempty"`
	Password string            `json:"password,omitempty"`
	ESIndex  map[string]string `json:"index"`
}
