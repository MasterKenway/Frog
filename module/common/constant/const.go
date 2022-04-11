package constant

import "time"

// Response Code Constant
const (
	CodeSuccess = "2000"

	CodeBadRequest       = "4003"
	CodeForbidden        = "4004"
	CodeLoginRequired    = "4001"
	CodePwdDecFailed     = "4002"
	CodePwdOrUsernameErr = "4003"
	CodeEmailCodeInvalid = "4004"
	CodeSendEmailFailed  = "4005"
	CodeCaptchaNeeded    = "4006"
	CodeCaptchaInvalid   = "4007"
	CodeUserExists       = "4008"
	CodeRateLimit        = "4009"

	CodeInternalError = "5000"
)

// Response Msg Constant
const (
	MsgSuccess = "Success"

	MsgApiNotExists      = "Api Not Exists"
	MsgIllegalRequest    = "Illegal Request"
	MsgParamInvalid      = "Request Param Invalid"
	MsgTimeStampOutdated = "Timestamp is Invalid"
	MsgNotLogin          = "Login Required"
	MsgPwdOrUsernameErr  = "Password or Username Error"
	MsgCaptchaNeeded     = "Captcha is Needed"
	MsgCaptchaInvalid    = "Captcha is Invalid"
	MsgEmailCodeInvalid  = "Email Code is Invalid"
	MsgUserExists        = "Username Exists"
	MsgRateLimit         = "request to frequently"

	MsgInternalError = "Internal Error"
)

// Header Constant
const (
	CookieKeyLoginCert = "x-login-status"

	HeaderKeyTimeStamp     = "x-timestamp"
	HeaderKeyNonce         = "x-nonce"
	HeaderKeyRequestID     = "x-request-id"
	HeaderKeyXForwardedFor = "x-forwarded-for"

	CtxKeyRemoteIP    = "remote-ip"
	CtxKeyRemotePort  = "remote-port"
	CtxKeyCmd         = "req-cmd"
	CtxKeyReqBody     = "req-body"
	CtxKeyRequestID   = "req-id"
	CtxKeyIsProxy     = "is-proxy"
	CtxKeyIsBot       = "is-bot"
	CtxKeyIsQuickConn = "is-quick-conn"
	CtxKeyUserInfo    = "x-user-info"
)

// RedisKey
const (
	RedisKeyNonce     = "req-nonce-"
	RedisKeyIPStamp   = "spider-ips-"
	RedisKeyEmailCode = "email-code-"
	RedisKeyLoginCert = "login-status-"
	RedisKeyRateLimit = "rate-limit-"
	RedisKeyApiCache  = "api-cache-"
)

// IP Stamp Constants
const (
	IPStampSpider    = 0
	IPStampBot       = 1
	IPStampProxy     = 2
	IPStampQuickConn = 3
)

// ETCD Key
const (
	// ETCD Config
	EtcdKeyMysqlConfig     = "mysql-config"
	EtcdKeyRedisConfig     = "redis-config"
	EtcdKeyKafkaConfig     = "kafka-config"
	EtcdKeyESConfig        = "elastic-config"
	EtcdKeyCaptchaConfig   = "captcha-config"
	EtcdKeyEmailConfig     = "email-config"
	EtcdKeyRateLimitConfig = "rate-limit-config"
)

// Kafka Key
const (
	KafkaKeyLogTopic = "es-log"
)

// Mysql Config Map Keys
const (
	MysqlUserInsertOnly = "insert_only"
	MysqlUserReadOnly   = "read_only"
	MysqlUserAll        = "all"
)

// URL constant
const (
	CaptchaDomain = "captcha.tencentcloudapi.com"
)

const (
	OneDay = 24 * time.Hour
)

const (
	FilterTypeEQ = "eq"
	FilterTypeGT = "gt"
	FilterTypeLT = "lt"
)
