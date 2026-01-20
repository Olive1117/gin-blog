package handler

import (
	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/service"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
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
	var registerRequest model.RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	var user model.User
	copier.CopyWithOption(&user, &registerRequest, copier.Option{IgnoreEmpty: true})
	if err := u.Service.Create(cx, &user); err != nil {
		errs.Fail(c, err)
		return
	}
	var userDTO model.UserDTO
	copier.CopyWithOption(&userDTO, &user, copier.Option{IgnoreEmpty: true})
	errs.Success(c, userDTO)
}
func (u *userHandler) Delete(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	if err := u.Service.Delete(cx, id); err != nil {
		errs.Fail(c, err)
	}
	errs.Success(c, nil)
}
func (u *userHandler) Get(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	user, err := u.Service.Get(cx, id)
	if err != nil {
		errs.Fail(c, err)
		return
	}
	var userDTO model.UserDTO
	copier.CopyWithOption(&userDTO, &user, copier.Option{IgnoreEmpty: true})
	errs.Success(c, userDTO)
}
func (u *userHandler) List(c *gin.Context) {
	cx := c.Request.Context()
	var filter model.User
	page := cast.ToInt(c.DefaultQuery("page", "1"))
	pageSize := cast.ToInt(c.DefaultQuery("page_size", "10"))
	if err := c.ShouldBindQuery(&filter); err != nil {
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	users, total, err := u.Service.List(cx, page, pageSize, &filter)
	if err != nil {
		errs.Fail(c, err)
		return
	}
	var userDTOs []model.UserDTO
	copier.CopyWithOption(&userDTOs, &users, copier.Option{IgnoreEmpty: true})
	errs.Success(c, gin.H{
		"list":  userDTOs,
		"total": total,
	})
}
func (u *userHandler) Update(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	if err := u.Service.Update(cx, &user, id); err != nil {
		errs.Fail(c, err)
		return
	}
	errs.Success(c, nil)
}
func (u *userHandler) GetMe(c *gin.Context) {
	cx := c.Request.Context()
	current_user, exists := c.Get("current_user")
	if !exists {
		errs.Fail(c, errs.ErrAuth)
		return
	}
	id, ok := current_user.(int64)
	if !ok {
		errs.Fail(c, errs.ErrAuthCheckTokenFail)
		return
	}
	user, err := u.Service.Get(cx, id)
	if err != nil {
		errs.Fail(c, err)
		return
	}
	var userDTO model.UserDTO
	copier.CopyWithOption(&userDTO, &user, copier.Option{IgnoreEmpty: true})
	errs.Success(c, userDTO)
}
