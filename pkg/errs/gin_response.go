package errs

import (
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/Olive1117/gin-blog/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success(ctx *gin.Context, data any) {
	if data == nil {
		data = gin.H{}
	}
	ctx.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "ok",
		Data: data,
	})
}

// Fail 统一处理错误返回
func Fail(ctx *gin.Context, err error) {
	var appErr *AppError
	// 尝试断言是否为自定义的 AppError
	if errors.As(err, &appErr) {
		ctx.JSON(appErr.HttpCode, Response{
			Code: appErr.Code,
			Msg:  appErr.Message,
			Data: gin.H{},
		})
		return
	}
	// 如果是未知错误，返回 500
	logger.FromContext(ctx).Error("服务器内部错误",
		zap.Error(err),
		zap.String("stack", string(debug.Stack())), // 生产环境可关闭
	)
	ctx.JSON(http.StatusInternalServerError, Response{
		Code: 500,
		Msg:  "服务器内部错误",
		Data: gin.H{},
	})
}
