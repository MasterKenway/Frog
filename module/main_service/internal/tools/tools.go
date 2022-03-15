package tools

import (
	"frog/module/common/constant"
	"frog/module/common/model/api_models"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

func CtxAbortWithCodeAndMsg(ctx *gin.Context, code string, msg string) {
	ctx.AbortWithStatusJSON(http.StatusOK, api_models.APIResponse{ResponseInfo: api_models.ResponseInfo{
		Code:    code,
		Message: msg,
	}})
	return
}

func RandInt() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(9481) + 6687
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func GetRedisKeyIPStamps(ip string) string {
	return constant.RedisKeyIPStamp + ip
}

func GetRedisKeyEmailCode(email string) string {
	return constant.RedisKeyEmailCode + email
}

func GetRedisKeyLoginCert(uid string) string {
	return constant.RedisKeyLoginCert + uid
}

func GetRedisKeyRateLimit(ip string) string {
	return constant.RedisKeyRateLimit + ip
}
