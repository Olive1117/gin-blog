package handler

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"github.com/gin-gonic/gin"
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

func (l *LoginHandler) Login(ctx *gin.Context) {
	var req model.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errs.Errors(ctx, errs.ErrInvalidParam)
		return
	}
	logger.FromContext(ctx.Request.Context()).Info("登录")
	res, err := l.Server.Login(ctx.Request.Context(), &req)
	if err != nil {
		errs.Errors(ctx, err)
		return
	}
	errs.Success(ctx, res)
}
