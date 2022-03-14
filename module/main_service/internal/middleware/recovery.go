package middleware

import (
	"frog/module/common/constant"
	"frog/module/main_service/internal/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"strings"
)

func GinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				//httpRequest, _ := httputil.DumpRequest(ctx.Request, false)
				if brokenPipe {
					log.Errorf(ctx.GetString(constant.CtxKeyRequestID), "err: %+v", err)
					//logger.Error(ctx.Request.URL.Path,
					//	zap.Any("error", err),
					//	zap.String("request", string(httpRequest)),
					//)
					// If the connection is dead, we can't write a status to it.
					ctx.Error(err.(error)) // nolint: errcheck
					ctx.Abort()
					return
				}

				if stack {
					//logger.Error("[Recovery from panic]",
					//	zap.Any("error", err),
					//	zap.String("request", string(httpRequest)),
					//	zap.String("stack", string(debug.Stack())),
					//)
				} else {
					//logger.Error("[Recovery from panic]",
					//	zap.Any("error", err),
					//	zap.String("request", string(httpRequest)),
					//)
				}
				ctx.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		ctx.Next()
	}
}
