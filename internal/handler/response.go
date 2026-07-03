package handler

import (
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Success(ctx *gin.Context, data any) {
	if data == nil {
		data = gin.H{}
	}
	ctx.JSON(http.StatusOK, model.Response{
		Code: 200,
		Msg:  "ok",
		Data: data,
	})
}

// Fail 统一处理错误返回
func Fail(ctx *gin.Context, err error) {
	var appErr *errs.AppError
	// 尝试断言是否为自定义的 AppError
	if errors.As(err, &appErr) {
		ctx.JSON(appErr.HttpCode, model.Response{
			Code: appErr.Code,
			Msg:  appErr.Message,
			Data: gin.H{},
		})
		return
	}
	// 如果是未知错误，返回 500
	logger.ErrorContext(ctx, "服务器内部错误",
		logger.Err(err),
		logger.String("stack", string(debug.Stack())), // 生产环境可关闭
	)
	ctx.JSON(http.StatusInternalServerError, model.Response{
		Code: 500,
		Msg:  "服务器内部错误",
		Data: gin.H{},
	})
}

func NewPageResponse[T any](list []T, total int64, q model.PageQuery) model.PageResponse[T] {
	return model.PageResponse[T]{
		List:     list,
		Page:     q.Page,
		PageSize: q.PageSize,
		Total:    total,
	}
}
