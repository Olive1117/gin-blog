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

type categoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(service service.CategoryService) CategoryHandler {
	return &categoryHandler{
		service: service,
	}
}

func (ch *categoryHandler) Create(c *gin.Context) {
	cx := c.Request.Context()
	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		Fail(c, errs.ErrInvalidParam)
		return
	}
	logger.DebugContext(cx, "创建分类", logger.Any("分类", category))

	if err := ch.service.Create(cx, &category); err != nil {
		Fail(c, err)
		return
	}
	categoryVO := convert.CategoryToVO(&category)
	Success(c, categoryVO)
}

func (ch *categoryHandler) List(c *gin.Context) {
	cx := c.Request.Context()
	pageQue := model.PageQuery{Page: cast.ToInt(c.DefaultQuery("page", "1")), PageSize: cast.ToInt(c.DefaultQuery("page_size", "10"))}
	var filter model.Category
	if err := c.ShouldBindQuery(&filter); err != nil {
		Fail(c, errs.ErrInvalidParam)
		return
	}
	logger.DebugContext(cx, "获取分类列表", logger.Any("过滤器", filter))

	pageRes, err := ch.service.List(cx, pageQue, &filter)
	if err != nil {
		Fail(c, err)
		return
	}
	categoryVOs := convert.MapSlice(pageRes.List, convert.CategoryToVO)
	Success(c, NewPageResponse(categoryVOs, pageRes.Total, pageQue))
}

func (ch *categoryHandler) Get(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	logger.DebugContext(cx, "获取分类详情", logger.Int64("分类ID", id))
	category, err := ch.service.Get(cx, id)
	if err != nil {
		Fail(c, err)
		return
	}
	categoryVO := convert.CategoryToVO(&category)
	Success(c, categoryVO)
}

func (ch *categoryHandler) Update(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		Fail(c, errs.ErrInvalidParam)
		return
	}
	logger.DebugContext(cx, "更新分类", logger.Any("分类", category))
	if err := ch.service.Update(cx, &category, id); err != nil {
		Fail(c, err)
		return
	}
	Success(c, nil)
}

func (ch *categoryHandler) Delete(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	logger.DebugContext(cx, "删除分类", logger.Int64("分类ID", id))
	err := ch.service.Delete(cx, id)
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c, nil)
}
