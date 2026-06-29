package handler

import (
	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/model/convert"
	"github.com/Olive1117/gin-blog/internal/service"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
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
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	logger.FromContext(cx).Debug("创建分类", zap.Any("分类", category))

	if err := ch.service.Create(cx, &category); err != nil {
		errs.Fail(c, err)
		return
	}
	categoryVO := convert.CategoryToVO(&category)
	errs.Success(c, categoryVO)
}

func (ch *categoryHandler) List(c *gin.Context) {
	cx := c.Request.Context()
	page := cast.ToInt(c.DefaultQuery("page", "1"))
	pageSize := cast.ToInt(c.DefaultQuery("page_size", "10"))
	var filter model.Category
	if err := c.ShouldBindQuery(&filter); err != nil {
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	logger.FromContext(cx).Debug("获取分类列表", zap.Any("过滤器", filter))

	categories, total, err := ch.service.List(cx, page, pageSize, &filter)
	if err != nil {
		errs.Fail(c, err)
		return
	}
	categoryVOs := convert.MapSlice(categories, convert.CategoryToVO)
	errs.Success(c, gin.H{
		"list":      categoryVOs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (ch *categoryHandler) Get(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	logger.FromContext(cx).Debug("获取分类详情", zap.Int64("分类ID", id))
	category, err := ch.service.Get(cx, id)
	if err != nil {
		errs.Fail(c, err)
		return
	}
	categoryVO := convert.CategoryToVO(&category)
	errs.Success(c, categoryVO)
}

func (ch *categoryHandler) Update(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	logger.FromContext(cx).Debug("更新分类", zap.Any("分类", category))
	if err := ch.service.Update(cx, &category, id); err != nil {
		errs.Fail(c, err)
		return
	}
	errs.Success(c, nil)
}

func (ch *categoryHandler) Delete(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	logger.FromContext(cx).Debug("删除分类", zap.Int64("分类ID", id))
	err := ch.service.Delete(cx, id)
	if err != nil {
		errs.Fail(c, err)
		return
	}
	errs.Success(c, nil)
}
