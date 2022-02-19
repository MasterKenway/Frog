package constant

// Response Code Constant
const (
	CodeSuccess = "200"

	CodeBadRequest = "403"
	CodeForbidden  = "404"

	CodeInternalError = "500"
)

// Response Msg Constant
const (
	MsgSuccess = "Success"

	MsgApiNotExists             = "Api Not Exists"
	MsgIllegalRequest           = "Illegal Request"
	MsgParamInvalid             = "Request Param Invalid"
	MsgTimeStampOutdated        = "Timestamp is Invalid"
	MsgTicketOrRandStrNotExists = "Ticket or RandStr Not Exists"
	MsgCaptchaValidateFailed    = "Captcha Validate Failed"

	MsgInternalError = "Internal Error"
)

// Header Constant
const (
	HeaderKeyTimeStamp = "x-timestamp"
	HeaderKeyNonce     = "x-nonce"
	HeaderKeyRequestID = "x-request-id"

	CtxKeyRemoteIP    = "remote-ip"
	CtxKeyRemotePort  = "remote-port"
	CtxKeyCmd         = "req-cmd"
	CtxKeyReqBody     = "req-body"
	CtxKeyRequestID   = "req-id"
	CtxKeyIsProxy     = "is-proxy"
	CtxKeyIsBot       = "is-bot"
	CtxKeyIsQuickConn = "is-quick-conn"
)

// RedisKey
const (
	RedisKeyNonce   = "req-nonce-"
	RedisKeyIPStamp = "spider-ips"
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
	EtcdKeyCaptchaConfig = "captcha-config"
)

// URL constant
const (
	CaptchaDomain = "captcha.tencentcloudapi.com"
)

//const (
//	CaptchaAppID        = "2031922961"
//	CaptchaAppSecretKey = "08Fa6z3Lod0sEOwPVSCo3zg**"
//)
