package errs

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": data,
	})
}

// HandleError 统一处理错误返回
func Errors(ctx *gin.Context, err error) {
	var appErr *AppError
	// 尝试断言是否为自定义的 AppError
	if errors.As(err, &appErr) {
		ctx.JSON(appErr.HttpCode, gin.H{
			"code": appErr.Code,
			"msg":  appErr.Message,
			"data": nil,
		})
		return
	}
	// 如果是未知错误，返回 500
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"code": 500,
		"msg":  "服务器内部错误",
		"data": nil,
	})
}
