package config

type MysqlConfig struct {
	Endpoint string `json:"endpoint,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
}
