package service

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/repository"
	"github.com/Olive1117/gin-blog/pkg/database"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"go.uber.org/zap"
)

type ArticleService struct {
	Repo *repository.ArticleRepo
	Ts   *database.GormTransaction
}

func (a *ArticleService) Update(c context.Context, id uint, art model.ArticleDTO) error {
	return a.Ts.Transaction(c, func(c context.Context) error {
		category, err := a.Repo.SyncCategory(c, art.Category)
		if err != nil {
			return err
		}
		tags, err := a.Repo.SyncTags(c, art.Tags)
		if err != nil {
			return err
		}
		article := &model.Article{
			Title:      art.Title,
			Desc:       art.Desc,
			Content:    art.Content,
			State:      art.State,
			CategoryID: category.ID,
			Tags:       tags,
		}
		article.ID = id
		logger.FromContext(c).Debug("更新文章业务", zap.Uint("id", id), zap.Any("文章", article))
		return a.Repo.Update(c, id, article)
	})
	// return a.Ts.Transaction(c, func(c context.Context) error {
	// 	if _, ok := artMap["Category"].(string); ok {

	// 	}
	// 	return a.Repo.Update(c, id, artMap)
	// })
}

func (a *ArticleService) Remove(c context.Context, id uint) (int, error) {
	rowsAffected, err := a.Repo.Delete(c, id)
	return rowsAffected, err
}

func NewArticleService(repo *repository.ArticleRepo, ts *database.GormTransaction) *ArticleService {
	return &ArticleService{
		Repo: repo,
		Ts:   ts,
	}
}
func (a *ArticleService) GetUserByID(c context.Context, id uint) (*model.Article, error) {
	article, err := a.Repo.GetById(c, id, "Category", "Tags")
	logger.FromContext(c).Debug("获取文章业务完成", zap.Uint("id", id), zap.Any("文章", article), zap.Error(err))
	// logger.FromContext(c).Debug("从数据库取出文章", zap.Any("文章", article))
	if err != nil {
		return nil, err
	}
	return &article, nil
}
func (a *ArticleService) Create(c context.Context, article *model.ArticleDTO) error {
	logger.FromContext(c).Debug("创建文章业务")
	return a.Ts.Transaction(c, func(c context.Context) error {
		// 1. 同步分类
		logger.FromContext(c).Debug("同步分类")
		category, err := a.Repo.SyncCategory(c, article.Category)
		if err != nil {
			return err
		}
		// 2. 同步标签
		logger.FromContext(c).Debug("同步标签")
		tags, err := a.Repo.SyncTags(c, article.Tags)
		if err != nil {
			return err
		}
		// 3. 构建实体
		article := &model.Article{
			Title:      article.Title,
			Desc:       article.Desc,
			Content:    article.Content,
			CategoryID: category.ID,
			Tags:       tags,
		}
		logger.FromContext(c).Debug("即将插入文章", zap.Any("文章", article))
		// 4. 调用 BaseRepo 提供的 Create 方法
		return a.Repo.CreateArticle(c, article)
	})
}
