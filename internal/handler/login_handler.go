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

func (l *LoginHandler) Login(gc *gin.Context) {
	c := gc.Request.Context()
	logger.FromContext(c).Debug("登录")
	var req model.LoginRequest
	if err := gc.ShouldBindJSON(&req); err != nil {
		logger.FromContext(c).Warn(errs.ErrInvalidParam.Message, zap.Any("err", err))
		errs.Fail(gc, errs.ErrInvalidParam)
		return
	}
	res, err := l.Server.Login(c, &req)
	if err != nil {
		logger.FromContext(c).Warn("[Business Warning]", zap.Any("err", err))
		errs.Fail(gc, err)
		return
	}
	errs.Success(gc, res)
}
