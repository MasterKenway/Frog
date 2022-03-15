package log

import (
	"bytes"
	"context"
	"frog/module/common/tools"
	"testing"
	"time"

	"frog/module/common/model/es_model"
	"frog/module/main_service/internal/config"
)

func TestESLogger(t *testing.T) {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Minute)
	)
	defer cancel()

	Info("test-log", "test")
	queryObject := []byte(`
{
    "query": {
        "bool": {
            "must": [
                {
                    "match": {
                        "request_id": "test-log"
                    }
                },
                {
                    "match": {
                        "message": "test"
                    }
                }
            ]
        }
    }
}
`)
	time.Sleep(15 * time.Second)

	resp, err := config.GetESCli().Search(
		config.GetESCli().Search.WithContext(ctx),
		config.GetESCli().Search.WithIndex(config.GetESIndexByConfig(es_model.ESLog{}.Index())),
		config.GetESCli().Search.WithBody(bytes.NewReader(queryObject)),
	)

	if err != nil {
		t.Fatalf("es search failed %s", err.Error())
	}

	if resp != nil {
		defer resp.Body.Close()

		if resp.IsError() {
			t.Fatalf("es search error %s", resp.String())
		}

		var hits []es_model.ESLog
		err := tools.GetModelFromESResp(resp, &hits)
		if err != nil {
			t.Fatalf("failed to get hits from resp, %s", err.Error())
		}
	}

}
