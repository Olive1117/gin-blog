package handler

import (
	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/model/convert"
	"github.com/Olive1117/gin-blog/internal/service"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
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
		Fail(c, errs.ErrInvalidParam)
		return
	}
	logger.DebugContext(cx, "创建标签", logger.Any("标签", tag))

	if err := th.service.Create(cx, &tag); err != nil {
		Fail(c, err)
		return
	}
	tagVO := convert.TagToVO(&tag)
	Success(c, tagVO)
}

func (th *tagHandler) List(c *gin.Context) {
	cx := c.Request.Context()
	var filter model.Tag
	pageQue := model.PageQuery{Page: cast.ToInt(c.DefaultQuery("page", "1")), PageSize: cast.ToInt(c.DefaultQuery("page_size", "10"))}
	if err := c.ShouldBindQuery(&filter); err != nil {
		Fail(c, errs.ErrInvalidParam)
		return
	}
	logger.DebugContext(cx, "获取标签列表", logger.Any("过滤器", filter))

	pageRes, err := th.service.List(cx, pageQue, &filter)
	if err != nil {
		Fail(c, err)
		return
	}
	tagVOs := convert.MapSlice(pageRes.List, convert.TagToVO)
	Success(c, NewPageResponse(tagVOs, pageRes.Total, pageQue))
}

func (th *tagHandler) Get(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	logger.DebugContext(cx, "获取标签", logger.Int64("标签ID", id))

	tag, err := th.service.Get(cx, id)
	if err != nil {
		Fail(c, err)
		return
	}
	tagVO := convert.TagToVO(&tag)
	Success(c, tagVO)
}

func (th *tagHandler) Delete(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	logger.DebugContext(cx, "删除标签", logger.Int64("标签ID", id))

	if err := th.service.Delete(cx, id); err != nil {
		Fail(c, err)
		return
	}
	Success(c, nil)
}

func (th *tagHandler) Update(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	var tag model.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		Fail(c, errs.ErrInvalidParam)
		return
	}
	logger.DebugContext(cx, "更新标签", logger.Any("标签", tag))
	if err := th.service.Update(cx, &tag, id); err != nil {
		Fail(c, err)
		return
	}
	Success(c, nil)
}
