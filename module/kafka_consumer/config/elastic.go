package config

import (
	"encoding/json"

	"frog/module/common/config"
	"frog/module/common/constant"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

var (
	esCli  *elasticsearch.Client
	esConf config.ElasticConfig
)

func GetESCli() *elasticsearch.Client {
	if esCli != nil {
		return esCli
	}

	esConfBytes, err := GetConfig(constant.EtcdKeyESConfig)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(esConfBytes, &esConf)
	if err != nil {
		panic(err)
	}

	esCli, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses:         esConf.Urls,
		Username:          esConf.Username,
		Password:          esConf.Password,
		MaxRetries:        5,
		EnableDebugLogger: true,
	})

	return esCli
}

func GetESIndexByConfig(configIndexKey string) string {
	return esConf.ESIndex[configIndexKey]
}

func GetESIndexer(index string) (esutil.BulkIndexer, error) {
	return esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  esConf.ESIndex[index],
		Client: GetESCli(),
	})
}
