package rental_info

import (
	"bytes"
	"context"
	"encoding/json"
	"frog/module/main_service/internal/config"
	"frog/module/main_service/internal/log"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"sync"
	"time"

	"frog/module/common/model/es_model"
	"frog/module/common/tools"
)

var (
	rentalInfoChannel = make(chan es_model.RentalInfo, 10000) // 接收创建帖子的结构体channel，批量写入到ES

	numMapMovement = map[int]float64{
		0: -171.2,
		1: -128.4,
		2: -192.6,
		3: -64.2,
		4: 0,
		5: -85.6,
		6: -149.8,
		7: -107,
		8: -42.8,
		9: -21.4,
	}
)

func init() {
	go ConsumeRentalInfo()
}

func ConsumeRentalInfo() {
	var (
		batchSize       = 1000
		consumerSize    = 1
		rentalInfoQueue = make([]es_model.RentalInfo, 0)
		lock            = sync.Mutex{}
		trigger         = make(chan struct{})
		ticker          = time.NewTicker(5 * time.Second)
	)

	for i := 0; i < consumerSize; i++ {
		go func() {
			for {
				select {
				case <-tools.Done:
					return
				case item := <-rentalInfoChannel:
					isFull := false
					lock.Lock()
					rentalInfoQueue = append(rentalInfoQueue, item)
					if len(rentalInfoQueue) > batchSize {
						isFull = true
					}
					lock.Unlock()

					if isFull {
						trigger <- struct{}{}
					}
				}
			}
		}()
	}

	for {
		select {
		case <-tools.Done:
			// avoid leak
			ticker.Stop()
			return
		case <-ticker.C:
			lock.Lock()
			esModels := rentalInfoQueue
			rentalInfoQueue = make([]es_model.RentalInfo, 0)
			lock.Unlock()
			consumeRentalInfos(esModels)
		case <-trigger:
			lock.Lock()
			esModels := rentalInfoQueue
			rentalInfoQueue = make([]es_model.RentalInfo, 0)
			lock.Unlock()
			consumeRentalInfos(esModels)
		}
	}

}

func consumeRentalInfos(esModels []es_model.RentalInfo) {
	var (
		ctx = context.Background()
	)

	indexer, err := config.GetESIndexer(es_model.RentalInfo{}.Index())
	if err != nil {
		log.Errorf("batch-to-es", "failed to get es indexer, %s", err.Error())
	}

	for _, model := range esModels {
		itemReqId := model.RequestId
		jsonData, _ := json.Marshal(model)
		err := indexer.Add(ctx, esutil.BulkIndexerItem{
			Action: "index",
			Body:   bytes.NewReader(jsonData),
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, respItem esutil.BulkIndexerResponseItem, err error) {
				log.Errorf(itemReqId, "bulk index failed, result: %s", respItem.Error.Reason)
				if err != nil {
					log.Errorf(itemReqId, "bulk index failed, err: %s", err.Error())
				}
			},
		})
		if err != nil {
			log.Errorf("batch-to-es", "indexer.Add %s", err.Error())
			continue
		}
	}

	err = indexer.Close(ctx)
	if err != nil {
		log.Errorf("batch-to-es", "indexer.Close %s", err.Error())
		return
	}

	stats := indexer.Stats()
	if stats.NumFailed > 0 {
		log.Errorf("batch-to-es", "stats.NumFailed %d", stats.NumFailed)
	}

}
