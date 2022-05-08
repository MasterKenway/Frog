package login

import (
	"encoding/json"
	"time"
)

type RedisUserInfo struct {
	Uid       string    `json:"Uid"`
	CookieKey string    `json:"CookieKey"`
	Username  string    `json:"Username,omitempty"`
	Email     string    `json:"Email,omitempty"`
	LoginIPs  []string  `json:"LoginIPs,omitempty"`
	UserAgent string    `json:"UserAgent,omitempty"`
	LoginTime time.Time `json:"LoginTime,omitempty"`
}

func (r RedisUserInfo) MarshalBinary() (data []byte, err error) {
	return json.Marshal(r)
}
