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

type authService struct {
	Repo repository.AuthRepo
	jwt  model.JWTHandler
}

func NewAuthService(store repository.AuthRepo, jwt model.JWTHandler) AuthService {
	return &authService{
		Repo: store,
		jwt:  jwt,
	}
}

func (l *authService) Auth(c context.Context, req *model.AuthRequest) (*model.AuthResponse, error) {
	logger.FromContext(c).Debug("登录业务代码")
	id, err := l.Repo.CheckAuth(c, req.Username, req.Password)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.FromContext(c).Warn(errs.ErrAuth.Message, zap.Error(err))
			return nil, errs.ErrAuth
		}
		return nil, errs.Error
	}
	token, expiresAt, err := l.jwt.GenerateToken(id, req.Username)
	if err != nil {
		logger.FromContext(c).Warn(errs.ErrAuthToken.Message, zap.Error(err))
		return nil, errs.ErrAuthToken
	}
	res := &model.AuthResponse{
		AccessToken: token,
		ExpiresAt:   expiresAt,
		TokenType:   "Bearer",
	}
	return res, nil
}
