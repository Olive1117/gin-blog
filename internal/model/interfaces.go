package model

import (
	"context"
	"time"

	"github.com/Olive1117/gin-blog/pkg/jwt"
)

type JWTHandler interface {
	GenerateToken(userID uint, username string) (string, time.Time, error)
	ParseToken(tokenString string) (*jwt.Claims, error)
}

type TransactionManager interface {
	Transaction(c context.Context, fn func(c context.Context) error) error
}
