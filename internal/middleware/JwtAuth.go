package middleware

import (
	"errors"
	"strings"

	"github.com/Olive1117/gin-blog/pkg/jwt"
	"github.com/Olive1117/gin-blog/pkg/response"
	"github.com/gin-gonic/gin"
)

func JwtAuth(j *jwt.JWTHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(ctx, response.ERROR_AUTH, nil)
			ctx.Abort()
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				response.Error(ctx, response.ERROR_AUTH_CHECK_TOKEN_TIMEOUT, nil)
			} else {
				response.Error(ctx, response.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
			}
			ctx.Abort()
			return
		}
		ctx.Set("current_user", claims.UserID)
		ctx.Next()
	}
}
