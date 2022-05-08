package config

import (
	"encoding/json"
	"fmt"
	"frog/module/common/constant"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"sync"
)

var (
	cosClient     *cos.Client
	cosConfig     *COSConfig
	cosOnce       = sync.Once{}
	cosConfigOnce = sync.Once{}
)

type COSConfig struct {
	Bucket       string `json:"bucket,omitempty"`
	AppID        string `json:"appid,omitempty"`
	ApiSecretID  string `json:"api_secret_id,omitempty"`
	ApiSecretKey string `json:"api_secret_key,omitempty"`
}

func GetCOSConfig() *COSConfig {
	if cosConfig != nil {
		return cosConfig
	}

	cosConfigOnce.Do(func() {
		confBytes, err := GetConfig(constant.EtcdKeyCosConfig)
		if err != nil {
			panic(err)
		}

		cosConfig = &COSConfig{}
		err = json.Unmarshal(confBytes, cosConfig)
		if err != nil {
			panic(err)
		}
	})

	return cosConfig
}

func GetCOSClient() *cos.Client {
	if cosClient != nil {
		return cosClient
	}

	cosOnce.Do(func() {
		confBytes, err := GetConfig(constant.EtcdKeyCosConfig)
		if err != nil {
			panic(err)
		}

		cosConfig = &COSConfig{}
		err = json.Unmarshal(confBytes, cosConfig)
		if err != nil {
			panic(err)
		}

		// 将 examplebucket-1250000000 和 COS_REGION 修改为用户真实的信息
		// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。https://console.cloud.tencent.com/cos5/bucket
		// COS_REGION 可以在控制台查看，https://console.cloud.tencent.com/cos5/bucket, 关于地域的详情见 https://cloud.tencent.com/document/product/436/6224
		u, _ := url.Parse(fmt.Sprintf("https://%s.cos.ap-guangzhou.myqcloud.com", cosConfig.Bucket))
		// 用于Get Service 查询，默认全地域 service.cos.myqcloud.com
		su, _ := url.Parse("https://service.cos.myqcloud.com")
		b := &cos.BaseURL{BucketURL: u, ServiceURL: su}
		// 1.永久密钥
		cosClient = cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  cosConfig.ApiSecretID,  // 替换为用户的 SecretId，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
				SecretKey: cosConfig.ApiSecretKey, // 替换为用户的 SecretKey，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
			},
		})
	})

	return cosClient
}
