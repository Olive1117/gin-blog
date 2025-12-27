package service

import (
	"github.com/Olive1117/gin-blog/model"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/jwt"
)

type loginStore interface {
	CheckLogin(string, string) (uint, error)
}

type loginService struct {
	Repository loginStore
	jwt        *jwt.JWTHandler
}

func NewLoginService(store loginStore, jwt *jwt.JWTHandler) *loginService {
	return &loginService{
		Repository: store,
		jwt:        jwt,
	}
}

func (l *loginService) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	id, err := l.Repository.CheckLogin(req.Username, req.Password)
	if err != nil {
		return nil, errs.ErrNotExistUser
	}
	token, expiresIn, err := l.jwt.GenerateToken(id, req.Username)
	if err != nil {
		return nil, errs.ErrAuthToken
	}
	res := &model.LoginResponse{
		AccessToken: token,
		ExpiresIn:   int64(expiresIn.Seconds()),
		TokenType:   "Bearer",
	}
	return res, nil
}
