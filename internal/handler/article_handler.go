package handler

import (
	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/service"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type articleHandler struct {
	service service.ArticleService
}

func NewArticleHandler(service service.ArticleService) ArticleHandler {
	return &articleHandler{
		service: service,
	}
}

func (a *articleHandler) Create(c *gin.Context) {
	cx := c.Request.Context()
	var articleDTO model.ArticleDTO
	var article model.Article
	if err := c.ShouldBindJSON(&articleDTO); err != nil {
		logger.FromContext(cx).Warn("参数错误", zap.Error(err))
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	logger.FromContext(cx).Debug("创建文章", zap.Any("文章", articleDTO))
	copier.CopyWithOption(&article, &articleDTO, copier.Option{IgnoreEmpty: true})

	if err := a.service.Create(cx, &article); err != nil {
		errs.Fail(c, err)
		return
	}
	errs.Success(c, articleDTO)
}
func (a *articleHandler) Delete(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))

	err := a.service.Delete(cx, id)
	if err != nil {
		errs.Fail(c, err)
		return
	}
	errs.Success(c, nil)
}
func (a *articleHandler) Get(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	logger.FromContext(cx).Debug("获取文章", zap.Int64("id", id))

	article, err := a.service.Get(cx, id)
	if err != nil {
		errs.Fail(c, err)
		return
	}
	var articleDTO model.ArticleDTO
	copier.CopyWithOption(&articleDTO, &article, copier.Option{IgnoreEmpty: true})
	errs.Success(c, articleDTO)
}
func (a *articleHandler) Update(c *gin.Context) {
	cx := c.Request.Context()
	var (
		articleDTO model.ArticleDTO
		article    model.Article
	)
	id := cast.ToInt64(c.Param("id"))
	err := c.ShouldBindJSON(&articleDTO)
	if err != nil {
		logger.FromContext(cx).Warn("参数错误", zap.Error(err))
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	logger.FromContext(cx).Debug("更新文章", zap.Int64("id", id), zap.Any("文章", articleDTO))
	copier.CopyWithOption(&article, &articleDTO, copier.Option{IgnoreEmpty: true})
	article.ID = id

	if err := a.service.Update(cx, &article); err != nil {
		errs.Fail(c, err)
		return
	}
	errs.Success(c, articleDTO)
}

func (a *articleHandler) List(c *gin.Context) {
	cx := c.Request.Context()
	var (
		articles    []model.Article
		articleDTOs []model.ArticleDTO
		article     model.Article
		query       model.ArticleQuery
	)
	page := cast.ToInt(c.DefaultQuery("page", "1"))
	pageSize := cast.ToInt(c.DefaultQuery("page_size", "10"))
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.FromContext(cx).Warn("参数错误", zap.Error(err))
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	copier.CopyWithOption(&article, &query, copier.Option{IgnoreEmpty: true})
	logger.FromContext(cx).Debug("获取文章列表", zap.Int("page", page), zap.Int("page_size", pageSize), zap.Any("query", &article))

	articles, total, err := a.service.List(cx, page, pageSize, &article)
	if err != nil {
		errs.Fail(c, err)
		return
	}
	copier.CopyWithOption(&articleDTOs, &articles, copier.Option{IgnoreEmpty: true})
	errs.Success(c, gin.H{
		"list":      articleDTOs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
