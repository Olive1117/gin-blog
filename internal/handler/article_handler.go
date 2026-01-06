package handler

import (
	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/service"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
)

type ArticleHandler struct {
	service *service.ArticleService
}

func NewArticleHandler(service *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		service: service,
	}
}

func (a *ArticleHandler) Create(c *gin.Context) {
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
func (a *ArticleHandler) Delete(c *gin.Context) {
	cx := c.Request.Context()
	id64, err := convertor.ToInt(c.Param("id"))
	var id uint = uint(id64)
	if err != nil {
		logger.FromContext(cx).Warn("参数错误", zap.Error(err))
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}

	rowsAffected, err := a.service.Delete(cx, id)
	if err != nil {
		errs.Fail(c, err)
		return
	}
	errs.Success(c, gin.H{"rowsAffected": rowsAffected})
}
func (a *ArticleHandler) Get(c *gin.Context) {
	cx := c.Request.Context()
	id64, err := convertor.ToInt(c.Param("id"))
	var id uint = uint(id64)
	logger.FromContext(cx).Debug("获取文章", zap.Uint("id", id))
	if err != nil {
		logger.FromContext(cx).Warn("参数错误", zap.Error(err))
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}

	article, err := a.service.Get(cx, id)
	if err != nil {
		errs.Fail(c, err)
		return
	}
	var articleDTO model.ArticleDTO
	copier.CopyWithOption(&articleDTO, &article, copier.Option{IgnoreEmpty: true})
	errs.Success(c, articleDTO)
}
func (a *ArticleHandler) Update(c *gin.Context) {
	var (
		id         uint
		articleDTO model.ArticleDTO
		article    model.Article
	)
	cx := c.Request.Context()
	id64, _ := convertor.ToInt(c.Param("id"))
	id = uint(id64)
	err := c.ShouldBindJSON(&articleDTO)
	if err != nil {
		logger.FromContext(cx).Warn("参数错误", zap.Error(err))
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	logger.FromContext(cx).Debug("更新文章", zap.Uint("id", id), zap.Any("文章", articleDTO))
	copier.CopyWithOption(&article, &articleDTO, copier.Option{IgnoreEmpty: true})
	article.ID = id

	if err := a.service.Update(cx, &article); err != nil {
		errs.Fail(c, err)
		return
	}
	errs.Success(c, articleDTO)
}

func (a *ArticleHandler) List(c *gin.Context) {
	var (
		articles    []model.Article
		articleDTOs []model.ArticleDTO
		article     model.Article
		query       model.ArticleQuery
	)
	cx := c.Request.Context()
	page64, _ := convertor.ToInt(c.DefaultQuery("page", "1"))
	pageSize64, _ := convertor.ToInt(c.DefaultQuery("page_size", "10"))
	page := int(page64)
	pageSize := int(pageSize64)
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.FromContext(cx).Warn("参数错误", zap.Error(err))
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	copier.CopyWithOption(&article, &query, copier.Option{IgnoreEmpty: true})
	logger.FromContext(cx).Debug("获取文章列表", zap.Int("page", page), zap.Int("page_size", pageSize), zap.Any("query", &article))

	articles, err := a.service.List(cx, page, pageSize, &article)
	if err != nil {
		errs.Fail(c, err)
		return
	}
	copier.CopyWithOption(&articleDTOs, &articles, copier.Option{IgnoreEmpty: true})
	res := map[string]any{
		"list":  articleDTOs,
		"total": len(articleDTOs),
	}
	errs.Success(c, &res)
}
