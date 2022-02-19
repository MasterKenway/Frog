package middleware

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"graduation-project/module/main_service/internal/config"
	"graduation-project/module/main_service/internal/constant"
	"graduation-project/module/main_service/internal/log"
	"graduation-project/module/main_service/internal/model/api_models"
	"graduation-project/module/main_service/internal/tools"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	perrors "github.com/pkg/errors"
)

var (
	captchaNeededUrl = map[string]bool{}
)

func Captcha(ctx *gin.Context) {
	var (
		reqId        = ctx.GetString(constant.CtxKeyRequestID)
		remoteIP     = ctx.GetString(constant.CtxKeyRemoteIP)
		reqBody, _   = ctx.Get(constant.CtxKeyReqBody)
		reqBodyBytes = reqBody.([]byte)

		url string
	)

	if _, ok := captchaNeededUrl[url]; ok {
		ticket, _ := jsonparser.GetString(reqBodyBytes, "Ticket")
		randStr, _ := jsonparser.GetString(reqBodyBytes, "RandStr")

		if ticket == "" || randStr == "" {
			tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeBadRequest, constant.MsgTicketOrRandStrNotExists)
			return
		}

		err := ValidateCaptcha(ticket, randStr, reqId, remoteIP)
		if err != nil {
			tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeBadRequest, constant.MsgCaptchaValidateFailed)
			return
		}
	}

	ctx.Next()
}

func ValidateCaptcha(ticket, randStr, reqId, ip string) error {
	//secretId, secretKey := config.GetCaptchaSecret()
	params := tools.TCParam{
		"Action":       "DescribeCaptchaResult",
		"Version":      "2019-07-22",
		"Nonce":        strconv.Itoa(tools.RandInt()),
		"Timestamp":    strconv.FormatInt(time.Now().Unix(), 10),
		"CaptchaType":  "9",
		"Ticket":       ticket,
		"Randstr":      randStr,
		"UserIp":       ip,
		"SecretId":     config.GetCaptchaConfig().ApiSecretID,
		"CaptchaAppId": config.GetCaptchaConfig().AppID,
		"AppSecretKey": config.GetCaptchaConfig().AppSecretKey,
	}

	params["Signature"] = params.CalHMACSHA1("GET", constant.CaptchaDomain, "/", config.GetCaptchaConfig().ApiSecretKey)

	resp, err := http.Get("https://" + constant.CaptchaDomain + "?" + params.GetUrlParam())
	if err != nil {
		log.Errorf("ValidateCaptcha: %s", err.Error())
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	var captchaResp api_models.CaptchaResponse
	err = json.Unmarshal(respBody, &captchaResp)
	if err != nil {
		log.Errorf(reqId, "failed unmarshal CaptchaResponse, err: %s", err.Error())
	}

	if captchaResp.RetCode != 0 && captchaResp.CaptchaResponseInfo.CaptchaCode != 1 {
		log.Errorf(reqId, "failed to verify, %+v", captchaResp)
		return perrors.Errorf("failed to verify, %+v", captchaResp)
	}

	return nil
}
