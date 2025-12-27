package handler

import (
	"github.com/Olive1117/gin-blog/model"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/gin-gonic/gin"
)

type loginStore interface {
	Login(*model.LoginRequest) (*model.LoginResponse, error)
}

type LoginHandler struct {
	Server loginStore
}

func NewLoginHandler(store loginStore) *LoginHandler {
	return &LoginHandler{Server: store}
}

func (l *LoginHandler) Login(ctx *gin.Context) {
	var req model.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errs.Errors(ctx, errs.ErrInvalidParam)
		return
	}
	res, err := l.Server.Login(&req)
	if err != nil {
		errs.Errors(ctx, err)
		return
	}
	errs.Success(ctx, res)
}
