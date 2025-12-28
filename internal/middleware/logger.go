package middleware

import (
	"time"

	"github.com/Olive1117/gin-blog/pkg/contextutil"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"github.com/duke-git/lancet/v2/random"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GinLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := ctx.GetHeader("X-Trace-ID")
		if traceID == "" {
			uuid, err := random.UUIdV4()
			if err != nil {
				logger.Log.Error("Trace生成错误")
			} else {
				traceID = uuid
			}
		}
		newctx := contextutil.SetTraceID(ctx, traceID)
		ctx.Request = ctx.Request.WithContext(newctx)
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		ctx.Next()
		cost := time.Since(start)
		logger.FromContext(ctx.Request.Context()).Info("HTTP Request",
			zap.Int("status", ctx.Writer.Status()),
			zap.String("method", ctx.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", ctx.ClientIP()),
			zap.Duration("cost", cost),
		)
	}
}
