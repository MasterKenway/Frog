package constant

// Response Code Constant
const (
	CodeSuccess      = "200"
	CodeParamInvalid = "403"
)

// Response Msg Constant
const (
	MsgSuccess           = "Success"
	MsgParamInvalid      = "Request Param Invalid"
	MsgTimeStampOutdated = "Timestamp is Invalid"
)

// Header Constant
const (
	HeaderKeyTimeStamp = "x-timestamp"
	HeaderKeyNonce     = "x-nonce"
	HeaderKeyRequestID = "x-request-id"
)
