package config

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"frog/module/common/tools"
	"frog/module/main_service/internal/log"

	perrors "github.com/pkg/errors"
	"go.etcd.io/etcd/client/v3"
)

var (
	etcdCli *clientv3.Client
)

// GetEtcdCli endpoint 错误会导致程序被阻塞
func GetEtcdCli() *clientv3.Client {
	if etcdCli != nil {
		return etcdCli
	}

	var err error
	etcdCli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"10.20.0.1:2379"},
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

func WatchConfig(key string, object interface{}) {
	if reflect.TypeOf(object).Kind() != reflect.Ptr {
		log.Errorf("etcd-watcher", "input param invalid")
		return
	}

	watchChan := GetEtcdCli().Watch(context.Background(), getEtcdKey(key))
	go func() {
		for {
			select {
			case <-tools.Done:
				return
			case resp := <-watchChan:
				if resp.Err() != nil {
					log.Errorf("etcd-watcher", "failed to watch key [%s], %s", getEtcdKey(key), resp.Err().Error())
					continue
				}

				for _, event := range resp.Events {
					if event.Type == clientv3.EventTypePut {
						err := json.Unmarshal(event.Kv.Value, object)
						if err != nil {
							log.Errorf("etcd-watcher", "failed to unmarshal object [%s], %s", getEtcdKey(key), err.Error())
						}
					}
				}
			}
		}

	}()
}

func getEtcdKey(key string) string {
	return tools.GetEnvType() + "-" + key
}
