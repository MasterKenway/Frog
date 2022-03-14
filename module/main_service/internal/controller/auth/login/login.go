package login

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"frog/module/main_service/internal/tools"
	"github.com/google/uuid"
	"time"

	"frog/module/common/constant"
	"frog/module/common/model/api_models"
	"frog/module/common/model/db_models"
	"frog/module/main_service/internal/config"
	"frog/module/main_service/internal/log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetLoginRequestAdapter() api_models.ApiInterface {
	return &Request{}
}

type Request struct {
	Username string `json:"username,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

func (r Request) GetResult(ctx *gin.Context) (interface{}, *api_models.APIError) {
	var (
		user  db_models.User
		reqId = ctx.GetString(constant.CtxKeyRequestID)
	)

	hash := md5.Sum([]byte(r.Password))
	err := config.GetReadOnlyMysqlCli().Model(&db_models.User{}).Where("username = ?", r.Username).First(&user).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Errorf(reqId, "failed to query userinfo from db, %s", err.Error())
			return nil, &api_models.APIError{
				Code:    constant.CodeInternalError,
				Message: constant.MsgInternalError,
			}
		} else {
			return nil, &api_models.APIError{
				Code:    constant.CodePwdOrUsernameErr,
				Message: constant.MsgPwdOrUsernameErr,
			}
		}
	}

	if user.Password != hex.EncodeToString(hash[:]) {
		return nil, &api_models.APIError{
			Code:    constant.CodePwdOrUsernameErr,
			Message: constant.MsgPwdOrUsernameErr,
		}
	}

	cookieKey := uuid.New().String()
	userInfo := RedisUserInfo{Username: user.Username, CookieKey: cookieKey}
	config.GetRedisCli().Set(context.Background(), tools.GetRedisKeyLoginCert(cookieKey), userInfo, 30*time.Minute)
	return userInfo, nil
}
