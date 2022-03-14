package es_model

import "time"

var (
	essLogMapping = `
{
    "mappings": {
        "properties": {
            "time": { "type": "date" },
            "level": { "type": "keyword" },
            "caller": { "type": "text" },
            "request_id": { "type": "keyword" },
            "message": { "type": "text" }
        }
    }
}
`
)

type ESLog struct {
	Time      time.Time `json:"time"`
	Level     string    `json:"level,omitempty"`
	Caller    string    `json:"caller,omitempty"`
	RequestID string    `json:"request_id,omitempty"`
	Message   string    `json:"message,omitempty"`
}

func (e ESLog) Index() string {
	return "es_log"
}

func (e ESLog) Mapping() string {
	return essLogMapping
}
