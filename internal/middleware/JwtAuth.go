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
			errs.Errors(ctx, errs.ErrLoginCheckTokenFail)
			ctx.Abort()
			return
		}
		parse := strings.SplitN(authHeader, " ", 2)
		if !(parse[0] == "Bearer" && len(parse) == 2) {
			errs.Errors(ctx, errs.ErrLoginCheckTokenFail)
			ctx.Abort()
			return
		}
		claims, err := j.ParseToken(parse[1])
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				errs.Errors(ctx, errs.ErrLoginCheckTokenTimeout)
			} else {
				errs.Errors(ctx, errs.ErrLoginCheckTokenFail)
			}
			ctx.Abort()
			return
		}
		ctx.Set("current_user", claims.UserID)
		ctx.Next()
	}
}
