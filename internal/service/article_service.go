package service

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/repository"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"go.uber.org/zap"
)

type articleService struct {
	Repo repository.ArticleRepo
	Ts   model.TransactionManager
}

func NewArticleService(repo repository.ArticleRepo, ts model.TransactionManager) ArticleService {
	return &articleService{
		Repo: repo,
		Ts:   ts,
	}
}

func (a *articleService) Update(c context.Context, article *model.Article) error {
	return a.Ts.Transaction(c, func(c context.Context) error {
		category, err := a.Repo.SyncCategory(c, article.Category.Name)
		if err != nil {
			return err
		}
		tagsName := make([]string, len(article.Tags))
		for i, tag := range article.Tags {
			tagsName[i] = tag.Name
		}
		tags, err := a.Repo.SyncTags(c, tagsName)
		if err != nil {
			return err
		}
		article.CategoryID = category.ID
		article.Category = *category
		article.Tags = tags
		logger.FromContext(c).Debug("更新文章业务", zap.Any("文章", article))
		return a.Repo.UpdateArticle(c, article)
	})
}

func (a *articleService) Create(c context.Context, article *model.Article) error {
	logger.FromContext(c).Debug("创建文章业务")
	return a.Ts.Transaction(c, func(c context.Context) error {
		// 1. 同步分类
		logger.FromContext(c).Debug("同步分类")
		category, err := a.Repo.SyncCategory(c, article.Category.Name)
		if err != nil {
			return err
		}
		// 2. 同步标签
		logger.FromContext(c).Debug("同步标签")
		tagsName := make([]string, len(article.Tags))
		for i, tag := range article.Tags {
			tagsName[i] = tag.Name
		}
		tags, err := a.Repo.SyncTags(c, tagsName)
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
		// 4. 调用 Create 方法
		return a.Repo.CreateArticle(c, article)
	})
}

func (a *articleService) Get(c context.Context, id uint) (model.Article, error) {
	return a.Repo.FindById(c, id, "Category", "Tags")
}

func (a *articleService) Delete(c context.Context, id uint) (int, error) {
	// 可以在这里增加删除后的逻辑（如清理 Redis）
	return a.Repo.Delete(c, id)
}

func (a *articleService) List(c context.Context, page, pageSize int, filter *model.Article) ([]model.Article, error) {
	return a.Repo.FindAllArticle(c, page, pageSize, filter)
}
