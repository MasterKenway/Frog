package config

import (
	"go.etcd.io/etcd/client/v3"
	"time"
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
		Endpoints:   []string{""},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	return etcdCli
}
