package service

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/repository"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/logger"
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
	logger.DebugContext(c, "登录业务代码")
	id, err := l.Repo.CheckAuth(c, req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	token, expiresAt, err := l.jwt.GenerateToken(id, req.Username)
	if err != nil {
		logger.WarnContext(c, errs.ErrAuthToken.Message, logger.Err(err))
		return nil, errs.ErrAuthToken
	}
	res := &model.AuthResponse{
		AccessToken: token,
		ExpiresAt:   expiresAt,
		TokenType:   "Bearer",
	}
	return res, nil
}
