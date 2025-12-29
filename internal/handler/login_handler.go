package handler

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LoginStore interface {
	Login(context.Context, *model.LoginRequest) (*model.LoginResponse, error)
}

type LoginHandler struct {
	Server LoginStore
}

func NewLoginHandler(store LoginStore) *LoginHandler {
	return &LoginHandler{Server: store}
}

func (l *LoginHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	logger.Debug(ctx, "登录")
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn(ctx, errs.ErrInvalidParam.Message, zap.Error(err))
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	res, err := l.Server.Login(ctx, &req)
	if err != nil {
		logger.Warn(ctx, "[Business Warning]", zap.Error(err))
		errs.Fail(c, err)
		return
	}
	errs.Success(c, res)
}
