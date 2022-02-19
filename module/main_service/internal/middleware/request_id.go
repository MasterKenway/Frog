package middleware

import (
	"graduation-project/module/main_service/internal/constant"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID(ctx *gin.Context) {
	ctx.Set(constant.CtxKeyRequestID, uuid.New().String())
}
