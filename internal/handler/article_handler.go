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
func (a *ArticleHandler) Delete(c *gin.Context) {}
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
	article, err := a.service.GetUserByID(cx, id)
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
func (a *ArticleHandler) Update(c *gin.Context) {}
func (a *ArticleHandler) List(c *gin.Context)   {}
