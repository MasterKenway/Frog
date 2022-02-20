package api_models

import (
	"encoding/json"
	"time"
)

type ControllerAdapter func() ApiInterface

type ApiInterface interface {
	GetResult() (interface{}, error)
}

type Request struct {
	Cmd  string          `json:"Cmd"`
	UUID string          `json:"UUID"`
	Data json.RawMessage `json:"Data"`
}

type APIResponse struct {
	ResponseInfo `json:"Response"`
}

type ResponseInfo struct {
	Code      string      `json:"Code"`
	Message   string      `json:"Message"`
	Data      interface{} `json:"Data"`
	Error     interface{} `json:"Error"`
	RequestID string      `json:"RequestID"`
}

// CaptchaResponse 与天御相应交互的结构体定义
type CaptchaResponse struct {
	CaptchaResponseInfo CaptchaResponseInfo `json:"Response"`
	RetCode             int64               `json:"retcode"`
	RetMsg              string              `json:"retmsg"`
}

type CaptchaResponseInfo struct {
	CaptchaCode    int    `json:"CaptchaCode"`
	CaptchaMsg     string `json:"CaptchaMsg"`
	EvilLevel      int    `json:"EvilLevel"`
	GetCaptchaTime int    `json:"GetCaptchaTime"`
	RequestId      string `json:"RequestId"`
}

type RawLog struct {
	Time      time.Time `json:"T"`
	Level     string    `json:"L,omitempty"`
	Caller    string    `json:"C,omitempty"`
	RequestID string    `json:"RID,omitempty"`
	Message   string    `json:"M,omitempty"`
}
