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
	if err := c.ShouldBindJSON(&articleDTO); err != nil {
		logger.WarnContext(cx, "参数错误", logger.Err(err))
		Fail(c, errs.ErrInvalidParam)
		return
	}
	logger.DebugContext(cx, "创建文章", logger.Any("文章", articleDTO))
	article := convert.ArticleFromDTO(&articleDTO)

	if err := a.service.Create(cx, article); err != nil {
		Fail(c, err)
		return
	}
	Success(c, articleDTO)
}
func (a *articleHandler) Delete(c *gin.Context) {
	cx := c.Request.Context()
	id := utils.ParseID(c.Param("id"))

	if err := a.service.Delete(cx, id); err != nil {
		Fail(c, err)
		return
	}
	Success(c, nil)
}
func (a *articleHandler) Get(c *gin.Context) {
	cx := c.Request.Context()

	id := utils.ParseID(c.Param("id"))
	logger.DebugContext(cx, "获取文章", logger.Int64("id", id))

	article, err := a.service.Get(cx, id)
	if err != nil {
		Fail(c, err)
		return
	}
	var articleVO = convert.ArticleToVO(&article)

	Success(c, articleVO)
}
func (a *articleHandler) Update(c *gin.Context) {
	cx := c.Request.Context()
	var (
		articleDTO model.ArticleDTO
		article    model.Article
	)
	id := utils.ParseID(c.Param("id"))
	logger.DebugContext(cx, "文章id", logger.Int64("id", id))
	err := c.ShouldBindJSON(&articleDTO)
	if err != nil {
		logger.WarnContext(cx, "参数错误", logger.Err(err))
		Fail(c, errs.ErrInvalidParam)
		return
	}
	logger.DebugContext(cx, "更新文章", logger.Int64("id", id), logger.Any("文章", articleDTO))
	article = *convert.ArticleFromDTO(&articleDTO)

	if err := a.service.Update(cx, &article, id); err != nil {
		Fail(c, err)
		return
	}
	Success(c, articleDTO)
}

func (a *articleHandler) List(c *gin.Context) {
	cx := c.Request.Context()
	var (
		article model.Article
		query   model.ArticleQuery
	)
	pageQue := model.PageQuery{Page: cast.ToInt(c.DefaultQuery("page", "1")), PageSize: cast.ToInt(c.DefaultQuery("page_size", "10"))}
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.WarnContext(cx, "参数错误", logger.Err(err))
		Fail(c, errs.ErrInvalidParam)
		return
	}
	article = *convert.ArticleFromQuery(&query)

	pageRes, err := a.service.List(cx, pageQue, &article)
	if err != nil {
		Fail(c, err)
		return
	}
	articleVOs := convert.MapSlice(pageRes.List, convert.ArticleToVO)
	Success(c, NewPageResponse(articleVOs, pageRes.Total, pageQue))
}

func (a *articleHandler) Stats(c *gin.Context) {
	cx := c.Request.Context()
	articleStats, err := a.service.Stats(cx)
	if err != nil {
		Fail(c, err)
		return
	}
	Success(c, articleStats)
}
