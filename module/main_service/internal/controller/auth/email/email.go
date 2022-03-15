package email

import (
	"context"
	"fmt"
	"time"

	"frog/module/common/constant"
	"frog/module/common/model/api_models"
	"frog/module/main_service/internal/config"
	"frog/module/main_service/internal/log"
	"frog/module/main_service/internal/tools"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func GetEmailCodeAdapter() api_models.ApiInterface {
	return &Request{}
}

type Request struct {
	Mail string `json:"mail" validate:"required,email"`
}

func (r *Request) GetResult(ctx *gin.Context) (interface{}, *api_models.APIError) {
	var (
		redisKey = tools.GetRedisKeyEmailCode(r.Mail)
	)

	randStr, err := config.GetRedisCli().Get(context.Background(), redisKey).Result()
	if err != redis.Nil {
		err = tools.SendEmail(r.Mail, fmt.Sprintf("验证码: %s (五分钟内有效)", randStr))
		if err != nil {
			return nil, &api_models.APIError{
				Code:    constant.CodeSendEmailFailed,
				Message: err.Error(),
			}
		}
		return nil, nil
	} else {
		log.Warnf("failed to get %s from redis, %s", redisKey, err.Error())
	}

	randStr = tools.RandStr(5)
	err = config.GetRedisCli().Set(context.Background(), redisKey, randStr, 5*time.Minute).Err()
	if err != nil {
		log.Errorf("failed to set %s:%s to redis, %s", redisKey, randStr, err.Error())
		return nil, &api_models.APIError{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}
	}

	err = tools.SendEmail(r.Mail, fmt.Sprintf("验证码: %s (五分钟内有效)", randStr))
	if err != nil {
		return nil, &api_models.APIError{
			Code:    constant.CodeSendEmailFailed,
			Message: err.Error(),
		}
	}

	return nil, nil
}
