package service

import (
	"context"
	"errors"

	"github.com/Olive1117/gin-blog/internal/handler"
	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/jwt"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var _ handler.LoginStore = (*LoginService)(nil)

type LoginStore interface {
	CheckLogin(context.Context, string, string) (uint, error)
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

func (l *LoginService) Login(c context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	logger.Debug(c, "登录业务代码")
	id, err := l.Repository.CheckLogin(c, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn(c, errs.ErrLogin.Message, zap.Error(err))
			return nil, errs.ErrLogin
		}
		return nil, errs.Error
	}
	token, expiresIn, err := l.jwt.GenerateToken(id, req.Username)
	if err != nil {
		logger.Warn(c, errs.ErrLoginToken.Message, zap.Error(err))
		return nil, errs.ErrLoginToken
	}
	res := &model.LoginResponse{
		AccessToken: token,
		ExpiresIn:   int64(expiresIn.Seconds()),
		TokenType:   "Bearer",
	}
	return res, nil
}
