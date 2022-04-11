package config

type MysqlConfig struct {
	ConfigMaps map[string]MysqlUserConfig `json:"config_maps"`
}

type MysqlUserConfig struct {
	Endpoint string `json:"endpoint,omitempty"`
	Port     int    `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	DBName   string `json:"db_name,omitempty"`
}
