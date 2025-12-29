package middleware

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
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
				logger.Log.Warn("Trace生成错误")
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
			zap.String("user-agent", ctx.Request.UserAgent()),
			zap.String("errors", ctx.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 检查断开的连接，因为它不是保证紧急堆栈跟踪的真正条件。
				var brokenPipe bool
				// OpError 是 net 包中的函数通常返回的错误类型。它描述了错误的操作、网络类型和地址。
				if ne, ok := err.(*net.OpError); ok {
					// SyscallError 记录来自特定系统调用的错误。
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") {
							brokenPipe = true
						}
					}
				}

				// DumpRequest 以 HTTP/1.x 连线形式返回给定的请求
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Log.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// 如果连接死了，我们就不能给它写状态
					c.Error(err.(error))
					c.Abort() // 终止该中间件
					return
				}

				if stack {
					logger.Log.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())), // 返回调用它的goroutine的格式化堆栈跟踪。
					)
				} else {
					logger.Log.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				// 调用 `Abort()` 并使用指定的状态代码写入标头。
				// StatusInternalServerError:500
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
