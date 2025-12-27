package middleware

import (
	"errors"
	"strings"

	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func JwtAuth(j *jwt.JWTHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			errs.Error(ctx, errs.ERROR_AUTH, nil)
			ctx.Abort()
			return
		}
		parse := strings.SplitN(authHeader, " ", 2)
		if !(parse[0] == "Bearer" && len(parse) == 2) {
			errs.Error(ctx, errs.ERROR_AUTH, nil)
			ctx.Abort()
			return
		}
		claims, err := j.ParseToken(parse[1])
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				errs.Error(ctx, errs.ERROR_AUTH_CHECK_TOKEN_TIMEOUT, nil)
			} else {
				errs.Error(ctx, errs.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
			}
			ctx.Abort()
			return
		}
		ctx.Set("current_user", claims.UserID)
		ctx.Next()
	}
}
