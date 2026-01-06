package handler

import (
	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/service"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gin-gonic/gin"
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
	var art model.ArticleDTO
	if err := c.ShouldBindJSON(&art); err != nil {
		logger.FromContext(cx).Warn("参数错误", zap.Error(err))
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	logger.FromContext(cx).Debug("创建文章", zap.Any("文章", art))
	if err := a.service.Create(cx, &art); err != nil {
		errs.Fail(c, err)
		return
	}
	errs.Success(c, art)
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
	rowsAffected, err := a.service.Repo.Delete(cx, id)
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
	article, err := a.service.Repo.FindById(c, id, "Category", "Tags")
	if err != nil {
		errs.Fail(c, err)
		return
	}
	res := &model.ArticleDTO{
		Title:    article.Title,
		Desc:     article.Desc,
		Content:  article.Content,
		State:    article.State,
		Category: article.Category.Name,
	}
	for _, tag := range article.Tags {
		res.Tags = append(res.Tags, tag.Name)
	}
	errs.Success(c, res)
}
func (a *ArticleHandler) Update(c *gin.Context) {
	var (
		id  uint
		art model.ArticleDTO
		err error
	)
	cx := c.Request.Context()
	id64, err := convertor.ToInt(c.Param("id"))
	id = uint(id64)
	err = c.ShouldBindJSON(&art)
	if err != nil {
		logger.FromContext(cx).Warn("参数错误", zap.Error(err))
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	logger.FromContext(cx).Debug("更新文章", zap.Uint("id", id), zap.Any("文章", art))
	err = a.service.Update(cx, id, art)
	errs.Success(c, art)
}

func (a *ArticleHandler) List(c *gin.Context) {
	cx := c.Request.Context()
	page64, _ := convertor.ToInt(c.DefaultQuery("page", "1"))
	pageSize64, _ := convertor.ToInt(c.DefaultQuery("page_size", "10"))
	page := int(page64)
	pageSize := int(pageSize64)
	query := &model.ArticleQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		logger.FromContext(cx).Warn("参数错误", zap.Error(err))
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	article := &model.Article{
		Title: query.Title,
		State: query.State,
	}
	article.Category.Name = query.Category
	for _, name := range query.Tags {
		article.Tags = append(article.Tags, model.Tag{Name: name})
	}
	logger.FromContext(cx).Debug("获取文章列表", zap.Int("page", page), zap.Int("page_size", pageSize), zap.Any("query", query))
	articles, err := a.service.Repo.FindAllArticle(cx, page, pageSize, article)
	if err != nil {
		errs.Fail(c, err)
		return
	}
	var res []model.ArticleDTO
	for _, article := range articles {
		item := model.ArticleDTO{
			Title:    article.Title,
			Desc:     article.Desc,
			Content:  article.Content,
			State:    article.State,
			Category: article.Category.Name,
		}
		for _, tag := range article.Tags {
			item.Tags = append(item.Tags, tag.Name)
		}
		res = append(res, item)
	}
	errs.Success(c, res)
}
