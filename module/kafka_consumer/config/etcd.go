package config

import (
	"context"
	"frog/module/common/tools"
	perrors "github.com/pkg/errors"
	"go.etcd.io/etcd/client/v3"
	"time"
)

var (
	etcdCli *clientv3.Client
)

func GetEtcdCli() *clientv3.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done := make(chan struct{})
	if etcdCli != nil {
		return etcdCli
	}

	go func() {
		var err error
		etcdCli, err = clientv3.New(clientv3.Config{
			Endpoints: []string{"10.10.0.4:2379"},
			//Endpoints:   []string{"10.20.0.1:2379"},
			DialTimeout: 2 * time.Second,
		})
		if err != nil {
			panic(err)
		}

		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		panic(perrors.New("connect with etcd timeout"))
	case <-done:
		return etcdCli
	}
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
	return tools.GetEnvType() + "-" + key
}
