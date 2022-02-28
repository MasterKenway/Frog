package config

type EmailConfig struct {
	EmailHost string `json:"host,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
}
