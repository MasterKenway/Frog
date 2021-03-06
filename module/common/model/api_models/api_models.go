package api_models

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
)

type ControllerAdapter func() ApiInterface

type ApiInterface interface {
	GetResult(ctx *gin.Context) (interface{}, *APIError)
}

type Request struct {
	Cmd  string          `json:"Cmd"`
	UUID string          `json:"UUID"`
	Data json.RawMessage `json:"Data"`
}

type Filter struct {
	Type  string      `json:"Type,omitempty" validate:"oneof=eq gt lt"`
	Key   string      `json:"Key,omitempty"`
	Value interface{} `json:"Value,omitempty"`
}

type APIResponse struct {
	ResponseInfo `json:"Response"`
}

type APIError struct {
	Code    string `json:"Code,omitempty"`
	Message string `json:"Message,omitempty"`
}

type ResponseInfo struct {
	Code      string      `json:"Code,omitempty"`
	Message   string      `json:"Message,omitempty"`
	Data      interface{} `json:"Data,omitempty"`
	Error     interface{} `json:"Error,omitempty"`
	RequestID string      `json:"RequestID,omitempty"`
}

// CaptchaResponse 与天御相应交互的结构体定义
type CaptchaResponse struct {
	CaptchaResponseInfo CaptchaResponseInfo `json:"Response"`
	RetCode             int64               `json:"retcode"`
	RetMsg              string              `json:"retmsg"`
}

type CaptchaRequest struct {
	RandStr string `json:"RandStr,omitempty"`
	Ticket  string `json:"Ticket,omitempty"`
}

type CaptchaResponseInfo struct {
	Error          interface{} `json:"Error"`
	CaptchaCode    int         `json:"CaptchaCode"`
	CaptchaMsg     string      `json:"CaptchaMsg"`
	EvilLevel      int         `json:"EvilLevel"`
	GetCaptchaTime int         `json:"GetCaptchaTime"`
	RequestId      string      `json:"RequestId"`
}

type RawLog struct {
	Time      time.Time `json:"T"`
	Level     string    `json:"L,omitempty"`
	Caller    string    `json:"C,omitempty"`
	RequestID string    `json:"RID,omitempty"`
	Message   string    `json:"M,omitempty"`
}
