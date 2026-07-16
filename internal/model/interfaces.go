package model

import (
	"context"
	"time"

	"github.com/Olive1117/gin-blog/pkg/jwt"
)

// jwt解码工具
type JWTHandler interface {
	GenerateAccessToken(userID string, roles []string) (string, time.Time, error)
	GenerateRefreshToken(userID string, roles []string, tokenID string) (string, time.Time, error)
	ParseToken(tokenString string) (*jwt.Claims, error)
}

// 事务工具
type TransactionManager interface {
	Transaction(c context.Context, fn func(c context.Context) error) error
}
