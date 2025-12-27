package handler

import (
	"errors"

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
		errs.Error(ctx, errs.INVALID_PARAMS, nil)
		return
	}
	res, err := l.Server.Login(&req)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrAuthToken):
			errs.Error(ctx, errs.ERROR_AUTH_TOKEN, nil)
		case errors.Is(err, errs.ErrNotExistUser):
			errs.Error(ctx, errs.ERROR_NOT_EXIST_USER, nil)
		default:
			errs.Error(ctx, errs.ERROR, nil)
		}
		return
	}
	errs.Success(ctx, res)
}
