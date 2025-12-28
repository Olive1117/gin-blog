package service

import (
	"context"
	"errors"

	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/jwt"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"gorm.io/gorm"
)

type LoginStore interface {
	CheckLogin(string, string) (uint, error)
}

type LoginService struct {
	Repository LoginStore
	jwt        *jwt.JWTHandler
}

func NewLoginService(store LoginStore, jwt *jwt.JWTHandler) *LoginService {
	return &LoginService{
		Repository: store,
		jwt:        jwt,
	}
}

func (l *LoginService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	logger.FromContext(ctx).Info("登录业务代码")
	id, err := l.Repository.CheckLogin(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrLogin
		}
		return nil, errs.Error
	}
	token, expiresIn, err := l.jwt.GenerateToken(id, req.Username)
	if err != nil {
		return nil, errs.ErrLoginToken
	}
	res := &model.LoginResponse{
		AccessToken: token,
		ExpiresIn:   int64(expiresIn.Seconds()),
		TokenType:   "Bearer",
	}
	return res, nil
}
