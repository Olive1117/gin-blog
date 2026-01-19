package handler

import (
	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/service"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type authHandler struct {
	Server service.AuthService
}

func NewAuthHandler(store service.AuthService) AuthHandler {
	return &authHandler{Server: store}
}

func (l *authHandler) Auth(c *gin.Context) {
	cx := c.Request.Context()
	logger.FromContext(cx).Debug("登录")
	var req model.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.FromContext(cx).Warn(errs.ErrInvalidParam.Message, zap.Error(err))
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	res, err := l.Server.Auth(cx, &req)
	if err != nil {
		logger.FromContext(cx).Warn("[Business Warning]", zap.Error(err))
		errs.Fail(c, err)
		return
	}
	errs.Success(c, res)
}
