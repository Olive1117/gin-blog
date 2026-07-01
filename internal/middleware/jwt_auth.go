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
)

func JwtAuth(j model.JWTHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		parse := strings.SplitN(c.GetHeader("Authorization"), " ", 2)
		if !(parse[0] == "Bearer" && len(parse) == 2) {
			logger.WarnContext(c.Request.Context(), errs.ErrAuthCheckTokenFail.Message)
			errs.Fail(c, errs.ErrAuthCheckTokenFail)
			c.Abort()
			return
		}
		logger.DebugContext(c.Request.Context(), "检查token", logger.String("token", parse[1]))
		claims, err := j.ParseToken(parse[1])
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				logger.WarnContext(c.Request.Context(), errs.ErrAuthCheckTokenTimeout.Message, logger.Err(err))
				errs.Fail(c, errs.ErrAuthCheckTokenTimeout)
			} else {
				logger.WarnContext(c.Request.Context(), errs.ErrAuthCheckTokenFail.Message, logger.Err(err))
				errs.Fail(c, errs.ErrAuthCheckTokenFail)
			}
			c.Abort()
			return
		}
		logger.DebugContext(c.Request.Context(), "注入当前用户id到上下文", logger.Int64("ID", claims.UserID))
		newctx := database.SetUserID(c.Request.Context(), claims.UserID)
		c.Request = c.Request.WithContext(newctx)
		c.Set("current_user", claims.UserID)
		c.Next()
	}
}
