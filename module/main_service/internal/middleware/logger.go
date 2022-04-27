package middleware

import (
	"time"

	"frog/module/common/constant"
	"frog/module/common/tools"
	"frog/module/main_service/internal/log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqId := uuid.New().String()
		ctx.Set(constant.CtxKeyRequestID, reqId)

		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		ctx.Next()
		cost := time.Since(start)
		log.Infof(reqId, "status: %d; method: %s; path %s; query: %s; ip: %s; user-agent: %s; errors: %s; cost: %d",
			ctx.Writer.Status(), ctx.Request.Method, path, query, tools.GetRemoteAddr(ctx),
			ctx.Request.UserAgent(), ctx.Errors.ByType(gin.ErrorTypePrivate).String(), cost)
	}
}
