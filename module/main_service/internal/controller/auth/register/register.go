package register

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"

	"frog/module/common/constant"
	"frog/module/common/model/api_models"
	"frog/module/common/model/db_models"
	"frog/module/main_service/internal/config"
	"frog/module/main_service/internal/log"
	"frog/module/main_service/internal/tools"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetRegisterAdapter() api_models.ApiInterface {
	return &Request{}
}

type Request struct {
	Username  string `json:"Username" validate:"required"`
	Password  string `json:"Password" validate:"required"`
	Email     string `json:"Email,omitempty" validate:"required"`
	EmailCode string `json:"EmailCode,omitempty" validate:"required"`
}

func (r *Request) GetResult(ctx *gin.Context) (interface{}, *api_models.APIError) {
	var (
		user              db_models.User
		reqId             = ctx.GetString(constant.CtxKeyRequestID)
		emailCodeRedisKey = tools.GetRedisKeyEmailCode(r.Email)
	)

	emailCode, err := config.GetRedisCli().Get(context.Background(), emailCodeRedisKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, &api_models.APIError{
				Code:    constant.CodeEmailCodeInvalid,
				Message: constant.MsgEmailCodeInvalid,
			}
		}

		log.Errorf(reqId, "failed to get %s from redis, %s", emailCodeRedisKey, err.Error())
		return nil, &api_models.APIError{
			Code:    constant.CodeInternalError,
			Message: constant.MsgInternalError,
		}
	}

	if emailCode != r.EmailCode {
		return nil, &api_models.APIError{
			Code:    constant.CodeEmailCodeInvalid,
			Message: constant.MsgEmailCodeInvalid,
		}
	}

	err = config.GetReadOnlyMysqlCli().Model(&db_models.User{}).Where("username = ?", r.Email).First(&user).Error
	if err != gorm.ErrRecordNotFound {
		if err != nil {
			log.Errorf(reqId, "failed to query user from db, %s", err.Error())
			return nil, &api_models.APIError{
				Code:    constant.CodeInternalError,
				Message: constant.MsgInternalError,
			}
		} else {
			return nil, &api_models.APIError{
				Code:    constant.CodeUserExists,
				Message: constant.MsgUserExists,
			}
		}
	}

	hash := md5.Sum([]byte(r.Password))
	user = db_models.User{
		Uid:      uuid.New().String(),
		Username: r.Username,
		Password: hex.EncodeToString(hash[:]),
		Email:    r.Email,
		LoginIps: sql.NullString{},
		IsValid:  0,
	}

	err = config.GetMysqlCli().Create(&user).Error
	if err != nil {
		log.Errorf(reqId, "failed to insert user into db, %s", err.Error())
		return nil, &api_models.APIError{
			Code:    constant.CodeInternalError,
			Message: constant.MsgInternalError,
		}
	}

	return nil, nil
}
