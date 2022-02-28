package login

import (
	"encoding/json"
	"time"
)

type RedisUserInfo struct {
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	LoginIPs  []string  `json:"login_ips,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	LoginTime time.Time `json:"login_time,omitempty"`
}

func (r RedisUserInfo) MarshalBinary() (data []byte, err error) {
	return json.Marshal(r)
}
