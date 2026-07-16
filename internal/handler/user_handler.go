package handler

import (
	"net/http"

	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/model/convert"
	"github.com/Olive1117/gin-blog/internal/service"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type userHandler struct {
	Service service.UserService
}

func NewUserHandler(service service.UserService) UserHandler {
	return &userHandler{
		Service: service,
	}
}
func (u *userHandler) Create(c *gin.Context) {
	cx := c.Request.Context()
	var UserDTO model.UserDTO
	if err := c.ShouldBindJSON(&UserDTO); err != nil {
		Fail(c, errs.ErrInvalidParam)
		return
	}
	user := convert.UserFromDTO(&UserDTO)
	if err := u.Service.Create(cx, user); err != nil {
		Fail(c, err)
		return
	}
	userVO := convert.UserToVO(user)
	Success(c, userVO)
}
func (u *userHandler) Delete(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	if err := u.Service.Delete(cx, id); err != nil {
		Fail(c, err)
	}
	Success(c, nil)
}
func (u *userHandler) Get(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	user, err := u.Service.Get(cx, id)
	if err != nil {
		Fail(c, err)
		return
	}
	logger.DebugContext(c.Request.Context(), "获取用户", logger.Any("用户", user))
	userVO := convert.UserToVO(&user)
	Success(c, userVO)
}
func (u *userHandler) List(c *gin.Context) {
	cx := c.Request.Context()
	var filter model.User
	pageQue := model.PageQuery{Page: cast.ToInt(c.DefaultQuery("page", "1")), PageSize: cast.ToInt(c.DefaultQuery("page_size", "10"))}
	if err := c.ShouldBindQuery(&filter); err != nil {
		Fail(c, errs.ErrInvalidParam)
		return
	}
	pageRes, err := u.Service.List(cx, pageQue, &filter)
	if err != nil {
		Fail(c, err)
		return
	}
	userVOs := convert.MapSlice(pageRes.List, convert.UserToVO)
	Success(c, NewPageResponse(userVOs, pageRes.Total, pageQue))
}
func (u *userHandler) Update(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		Fail(c, errs.ErrInvalidParam)
		return
	}
	if err := u.Service.Update(cx, &user, id); err != nil {
		Fail(c, err)
		return
	}
	Success(c, nil)
}
func (u *userHandler) GetMe(c *gin.Context) {
	cx := c.Request.Context()
	current_user, exists := c.Get("current_user")
	if !exists {
		Fail(c, errs.ErrAuth)
		return
	}
	id, ok := current_user.(int64)
	if !ok {
		Fail(c, errs.ErrAuthCheckTokenFail)
		return
	}
	user, err := u.Service.Get(cx, id)
	if err != nil {
		Fail(c, err)
		return
	}
	userVO := convert.UserToVO(&user)
	Success(c, userVO)
}
func (u *userHandler) Login(c *gin.Context) {
	cx := c.Request.Context()
	logger.DebugContext(cx, "登录")
	var userInfo = model.RefreshTokens{
		Device_info: c.GetHeader("X-Forwarded-For"),
		Ip_address:  c.GetHeader("X-Real-Ip"),
	}
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnContext(cx, errs.ErrInvalidParam.Message, logger.Err(err))
		Fail(c, errs.ErrInvalidParam)
		return
	}
	res, refreshCookie, err := u.Service.Login(cx, req, userInfo)
	if err != nil {
		logger.WarnContext(cx, "登录失败", logger.Err(err))
		Fail(c, err)
		return
	}
	c.SetCookieData(&http.Cookie{
		Name:  "__refresh_token",
		Value: refreshCookie.RefreshToken,
		// Path:     "/",                     // 限制路径，只允许刷新接口访问
		Expires:  refreshCookie.ExpiresAt, // 7天
		HttpOnly: true,                    // 禁止JS访问（防XSS）
		Secure:   true,                    // 仅HTTPS（防中间人攻击）
		SameSite: http.SameSiteLaxMode,    // 严格模式，仅同站请求和外站导航携带
	})
	Success(c, res)
}
func (u *userHandler) ChangePassword(c *gin.Context) {
	cx := c.Request.Context()
	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WarnContext(cx, "修改密码参数无效", logger.Err(err))
		Fail(c, errs.ErrInvalidParam)
		return
	}
	err := u.Service.ChangePassword(cx, req.Username, req.OldPassword, req.NewPassword)
	if err != nil {
		logger.WarnContext(cx, "修改密码失败", logger.Err(err))
		Fail(c, err)
		return
	}
	Success(c, req.NewPassword)
}
func (u *userHandler) RefreshToken(c *gin.Context) {
	cx := c.Request.Context()
	refreshToken, err := c.Cookie("__refresh_token")
	if err != nil {
		logger.WarnContext(cx, "获取刷新token失败", logger.Err(err))
		Fail(c, errs.ErrNotExistRefreshToken)
		return
	}
	res, err := u.Service.RefreshToken(cx, refreshToken)
	if err != nil {
		logger.WarnContext(cx, "获取刷新token失败", logger.Err(err))
		Fail(c, err)
		return
	}
	Success(c, res)
}
