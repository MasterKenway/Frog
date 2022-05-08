package tools

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	perrors "github.com/pkg/errors"
	"os"
	"reflect"
	"strings"

	"frog/module/common/constant"

	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
)

var (
	Done = make(chan struct{})

	EnvType string

	CronService = gocron.NewScheduler()
)

func GetEnvType() string {
	envType, ok := os.LookupEnv("EnvType")
	if !ok || envType == "" {
		EnvType = "Testing"
	} else {
		EnvType = envType
	}

	return EnvType
}

func GetRemoteAddr(ctx *gin.Context) string {
	var ip string
	ip = ctx.GetHeader(constant.HeaderKeyXForwardedFor)
	if ip != "" {
		return ip
	}
	return ctx.ClientIP()
}

func GetResultFromRedis(redisCli *redis.Client, key string, object interface{}) error {
	cacheBytes, err := redisCli.Get(context.Background(), key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(cacheBytes, object)
}

func GetModelCols(model interface{}) map[string]bool {
	colsMap := make(map[string]bool)
	modelType := reflect.TypeOf(model)
	for i := 0; i < modelType.NumField(); i++ {
		tags := modelType.Field(i).Tag.Get("gorm")
		parts := strings.Split(tags, ";")
		if len(parts) <= 0 {
			panic(perrors.New("error tags"))
		}
		kv := strings.Split(parts[0], ":")
		if len(kv) != 2 {
			panic(perrors.New("error tags"))
		}
		colsMap[kv[1]] = true
	}

	return colsMap
}

func Pagination(slice interface{}, pageNum, pageSize int) (reflect.Value, error) {
	sliceType := reflect.TypeOf(slice)
	res := reflect.MakeSlice(sliceType, 0, pageSize)
	if sliceType.Kind() == reflect.Slice {
		values := reflect.ValueOf(slice)
		for i := (pageNum - 1) * pageSize; i < values.Len() && i < pageNum*pageSize; i++ {
			res = reflect.Append(res, values.Index(i))
		}
	} else {
		return reflect.Value{}, perrors.New("type error")
	}
	return res, nil
}
