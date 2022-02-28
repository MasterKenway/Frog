package config

import (
	"context"
	"time"

	"frog/module/common"

	perrors "github.com/pkg/errors"
	"go.etcd.io/etcd/client/v3"
)

var (
	etcdCli *clientv3.Client
)

func GetEtcdCli() *clientv3.Client {
	if etcdCli != nil {
		return etcdCli
	}

	var err error
	etcdCli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://10.10.0.2:2379"},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	return etcdCli
}

func GetConfig(key string) ([]byte, error) {
	resp, err := GetEtcdCli().Get(context.Background(), getEtcdKey(key))
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) <= 0 {
		return nil, perrors.New("resp.Kvs nil")
	}

	return resp.Kvs[0].Value, nil
}

func getEtcdKey(key string) string {
	return common.GetEnvType() + "-" + key
}
