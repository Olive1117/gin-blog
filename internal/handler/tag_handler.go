package handler

import (
	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/service"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type tagHandler struct {
	service service.TagService
}

func NewTagHandler(service service.TagService) TagHandler {
	return &tagHandler{
		service: service,
	}
}

func (th *tagHandler) Create(c *gin.Context) {
	cx := c.Request.Context()
	var tag model.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	logger.FromContext(cx).Debug("创建标签", zap.Any("标签", tag))

	if err := th.service.Create(cx, &tag); err != nil {
		errs.Fail(c, err)
		return
	}
	errs.Success(c, tag)
}

func (th *tagHandler) List(c *gin.Context) {
	cx := c.Request.Context()
	var filter model.Tag
	page := cast.ToInt(c.DefaultQuery("page", "1"))
	pageSize := cast.ToInt(c.DefaultQuery("page_size", "10"))
	if err := c.ShouldBindQuery(&filter); err != nil {
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	logger.FromContext(cx).Debug("获取标签列表", zap.Any("过滤器", filter))

	tags, err := th.service.List(cx, page, pageSize, &filter)
	if err != nil {
		errs.Fail(c, err)
		return
	}
	errs.Success(c, tags)
}

func (th *tagHandler) Get(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	logger.FromContext(cx).Debug("获取标签", zap.Int64("标签ID", id))

	tag, err := th.service.Get(cx, id)
	if err != nil {
		errs.Fail(c, err)
		return
	}
	errs.Success(c, tag)
}

func (th *tagHandler) Delete(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	logger.FromContext(cx).Debug("删除标签", zap.Int64("标签ID", id))

	rowsAffected, err := th.service.Delete(cx, id)
	if err != nil {
		errs.Fail(c, err)
		return
	}
	errs.Success(c, gin.H{"rowsAffected": rowsAffected})
}

func (th *tagHandler) Update(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	var tag model.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	tag.ID = id
	logger.FromContext(cx).Debug("更新标签", zap.Any("标签", tag))
	if err := th.service.Update(cx, &tag); err != nil {
		errs.Fail(c, err)
		return
	}
	errs.Success(c, nil)
}
