package tools

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	perrors "github.com/pkg/errors"
	"io/ioutil"
	"reflect"
)

var (
	ErrESRecordNotFound = perrors.New("hits is empty")
)

type ESModel interface {
	Index() string
	Mapping() string
}

type ESResponse struct {
	Took    int
	Timeout bool
	Shard   struct {
		Total      int
		Successful int
		Skipped    int
		Failed     int
	}

	Hits struct {
		Total struct {
			Value    int
			Relation string
		}

		Hits json.RawMessage
	}
}

func GetModelFromESResp(resp *esapi.Response, target interface{}) error {
	if reflect.TypeOf(target).Kind() != reflect.Ptr {
		return perrors.Errorf("kind %s is invalid", reflect.TypeOf(target).String())
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	respObject := &ESResponse{}
	err = json.Unmarshal(respBytes, &respObject)
	if err != nil {
		return err
	}

	if respObject.Hits.Total.Value == 0 {
		return ErrESRecordNotFound
	}

	return json.Unmarshal(respObject.Hits.Hits, target)
}
