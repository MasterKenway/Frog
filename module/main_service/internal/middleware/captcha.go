package middleware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"frog/module/common/constant"
	"frog/module/common/model/api_models"
	"frog/module/main_service/internal/config"
	"frog/module/main_service/internal/log"
	"frog/module/main_service/internal/tools"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	perrors "github.com/pkg/errors"
)

var (
	captchaNeededUrl = map[string]bool{
		"CreateRegister": true,
	}
)

func Captcha() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			//reqId        = ctx.GetString(constant.CtxKeyRequestID)
			//remoteIP     = ctx.GetString(constant.CtxKeyRemoteIP)
			reqBody, _ = ctx.Get(constant.CtxKeyReqBody)
			cmd        = ctx.GetString(constant.CtxKeyCmd)

			reqBodyBytes = reqBody.([]byte)
		)

		if _, ok := captchaNeededUrl[cmd]; ok {
			ticket, _ := jsonparser.GetString(reqBodyBytes, "Ticket")
			randStr, _ := jsonparser.GetString(reqBodyBytes, "RandStr")

			if ticket == "" || randStr == "" {
				tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeCaptchaNeeded, constant.MsgCaptchaNeeded)
				return
			}

			//err := ValidateCaptcha(ticket, randStr, reqId, remoteIP)
			//if err != nil {
			//	tools.CtxAbortWithCodeAndMsg(ctx, constant.CodeCaptchaInvalid, constant.MsgCaptchaInvalid)
			//	return
			//}
		}

		ctx.Next()
	}
}

// ValidateCaptcha 这个勾八东西要钱，字符串不为空就直接放行
func ValidateCaptcha(ticket, randStr, reqId, ip string) error {
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
		return err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	fmt.Printf("%s", string(respBody))

	var captchaResp api_models.CaptchaResponse
	err = json.Unmarshal(respBody, &captchaResp)
	if err != nil {
		log.Errorf(reqId, "failed unmarshal CaptchaResponse, err: %s", err.Error())
		return err
	}

	if captchaResp.RetCode != 0 && captchaResp.CaptchaResponseInfo.CaptchaCode != 1 {
		log.Errorf(reqId, "failed to verify, %+v", captchaResp)
		return perrors.Errorf("failed to verify, %+v", captchaResp)
	}

	return nil
}
