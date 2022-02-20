package config

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"graduation-project/module/common/config"
	"graduation-project/module/common/constant"

	"github.com/olivere/elastic/v7"
)

var (
	esCli  *elastic.Client
	esConf config.ElasticConfig
)

func GetESCli() *elastic.Client {
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

	esCli, err = elastic.NewClient(
		elastic.SetURL(esConf.Urls...),
		elastic.SetBasicAuth(esConf.Username, esConf.Password),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(1*time.Minute),
		elastic.SetGzip(true),
		elastic.SetErrorLog(log.New(os.Stderr, "Elastic ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "Elastic ", log.LstdFlags)),
	)

	return esCli
}

func GetESIndexName(index string) string {
	return esConf.ESIndex[index]
}
