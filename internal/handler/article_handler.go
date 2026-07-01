package handler

import (
	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/model/convert"
	"github.com/Olive1117/gin-blog/internal/service"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"github.com/Olive1117/gin-blog/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
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
		logger.WarnContext(cx, "参数错误", logger.Err(err))
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	logger.DebugContext(cx, "创建文章", logger.Any("文章", articleDTO))
	article = *convert.ArticleFromDTO(&articleDTO)

	if err := a.service.Create(cx, &article); err != nil {
		errs.Fail(c, err)
		return
	}
	errs.Success(c, articleDTO)
}
func (a *articleHandler) Delete(c *gin.Context) {
	cx := c.Request.Context()
	var id int64
	param := c.Param("id")
	if utils.IsShortID(param) {
		id = utils.DecodeByOBID(param)
	} else {
		id = cast.ToInt64(param)
	}

	err := a.service.Delete(cx, id)
	if err != nil {
		errs.Fail(c, err)
		return
	}
	errs.Success(c, nil)
}
func (a *articleHandler) Get(c *gin.Context) {
	cx := c.Request.Context()
	var id int64
	param := c.Param("id")
	if utils.IsShortID(param) {
		id = utils.DecodeByOBID(param)
	} else {
		id = cast.ToInt64(param)
	}
	logger.DebugContext(cx, "获取文章", logger.Int64("id", id))

	article, err := a.service.Get(cx, id)
	if err != nil {
		errs.Fail(c, err)
		return
	}
	var articleVO = convert.ArticleToVO(&article)

	errs.Success(c, articleVO)
}
func (a *articleHandler) Update(c *gin.Context) {
	cx := c.Request.Context()
	var (
		articleDTO model.ArticleDTO
		article    model.Article
	)
	id := cast.ToInt64(c.Param("id"))
	logger.DebugContext(cx, "文章id", logger.Int64("id", id))
	err := c.ShouldBindJSON(&articleDTO)
	if err != nil {
		logger.WarnContext(cx, "参数错误", logger.Err(err))
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	logger.DebugContext(cx, "更新文章", logger.Int64("id", id), logger.Any("文章", articleDTO))
	article = *convert.ArticleFromDTO(&articleDTO)

	if err := a.service.Update(cx, &article, id); err != nil {
		errs.Fail(c, err)
		return
	}
	errs.Success(c, articleDTO)
}

func (a *articleHandler) List(c *gin.Context) {
	cx := c.Request.Context()
	var (
		articles   []model.Article
		articleVOs []model.ArticleVO
		article    model.Article
		query      model.ArticleQuery
	)
	page := cast.ToInt(c.DefaultQuery("page", "1"))
	pageSize := cast.ToInt(c.DefaultQuery("page_size", "10"))
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.WarnContext(cx, "参数错误", logger.Err(err))
		errs.Fail(c, errs.ErrInvalidParam)
		return
	}
	article = *convert.ArticleFromQuery(&query)
	logger.DebugContext(cx, "获取文章列表", logger.Int("page", page), logger.Int("page_size", pageSize), logger.Any("query", &article))

	articles, total, err := a.service.List(cx, page, pageSize, &article)
	if err != nil {
		errs.Fail(c, err)
		return
	}
	articleVOs = convert.MapSlice(articles, convert.ArticleToVO)
	errs.Success(c, gin.H{
		"list":      articleVOs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (a *articleHandler) Stats(c *gin.Context) {
	cx := c.Request.Context()
	articleStats, err := a.service.Stats(cx)
	if err != nil {
		errs.Fail(c, err)
		return
	}
	errs.Success(c, articleStats)
}
