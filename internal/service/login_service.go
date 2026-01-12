package service

import (
	"context"
	"errors"

	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/repository"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type loginService struct {
	Repo repository.LoginRepo
	jwt  model.JWTHandler
}

func NewLoginService(store repository.LoginRepo, jwt model.JWTHandler) LoginService {
	return &loginService{
		Repo: store,
		jwt:  jwt,
	}
}

func (l *loginService) Login(c context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	logger.FromContext(c).Debug("登录业务代码")
	id, err := l.Repo.CheckLogin(c, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.FromContext(c).Warn(errs.ErrLogin.Message, zap.Error(err))
			return nil, errs.ErrLogin
		}
		return nil, errs.Error
	}
	token, expiresAt, err := l.jwt.GenerateToken(id, req.Username)
	if err != nil {
		logger.FromContext(c).Warn(errs.ErrLoginToken.Message, zap.Error(err))
		return nil, errs.ErrLoginToken
	}
	res := &model.LoginResponse{
		AccessToken: token,
		ExpiresAt:   expiresAt,
		TokenType:   "Bearer",
	}
	return res, nil
}
