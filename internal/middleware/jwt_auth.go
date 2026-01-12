package middleware

import (
	"errors"
	"strings"

	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/pkg/database"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/jwt"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func JwtAuth(j model.JWTHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		// authHeader := c.GetHeader("Authorization")
		// if authHeader == "" {
		// 	logger.FromContext(c.Request.Context()).Warn(errs.ErrLoginCheckTokenFail.Message)
		// 	errs.Fail(c, errs.ErrLoginCheckTokenFail)
		// 	c.Abort()
		// 	return
		// }
		// parse := strings.SplitN(authHeader, " ", 2)
		parse := strings.SplitN(c.GetHeader("Authorization"), " ", 2)
		if !(parse[0] == "Bearer" && len(parse) == 2) {
			logger.FromContext(c.Request.Context()).Warn(errs.ErrLoginCheckTokenFail.Message)
			errs.Fail(c, errs.ErrLoginCheckTokenFail)
			c.Abort()
			return
		}
		logger.FromContext(c.Request.Context()).Debug("检查token", zap.String("token", parse[1]))
		claims, err := j.ParseToken(parse[1])
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				logger.FromContext(c.Request.Context()).Warn(errs.ErrLoginCheckTokenTimeout.Message, zap.Error(err))
				errs.Fail(c, errs.ErrLoginCheckTokenTimeout)
			} else {
				logger.FromContext(c.Request.Context()).Warn(errs.ErrLoginCheckTokenFail.Message, zap.Error(err))
				errs.Fail(c, errs.ErrLoginCheckTokenFail)
			}
			c.Abort()
			return
		}
		logger.FromContext(c.Request.Context()).Debug("注入当前用户id到上下文", zap.Uint("ID", claims.UserID))
		newctx := database.SetUserID(c.Request.Context(), claims.UserID)
		c.Request = c.Request.WithContext(newctx)
		c.Set("current_user", claims.UserID)
		c.Next()
	}
}
