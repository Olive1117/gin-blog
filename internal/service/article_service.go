package service

import (
	"context"
	"unicode/utf8"

	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/repository"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"go.uber.org/zap"
)

type articleService struct {
	Repo         repository.ArticleRepo
	Ts           model.TransactionManager
	TagRepo      repository.TagRepo
	CategoryRepo repository.CategoryRepo
}

func NewArticleService(repo repository.ArticleRepo, ts model.TransactionManager, tagRepo repository.TagRepo, categoryRepo repository.CategoryRepo) ArticleService {
	return &articleService{
		Repo:         repo,
		Ts:           ts,
		TagRepo:      tagRepo,
		CategoryRepo: categoryRepo,
	}
}

func (a *articleService) Update(c context.Context, article *model.Article, id int64) error {
	return a.Ts.Transaction(c, func(c context.Context) error {
		category, err := a.CategoryRepo.SyncCategory(c, article.Category.Name)
		if err != nil {
			return err
		}
		tagsName := make([]string, len(article.Tags))
		for i, tag := range article.Tags {
			tagsName[i] = tag.Name
		}
		tags, err := a.TagRepo.SyncTags(c, tagsName)
		if err != nil {
			return err
		}
		article.ID = id
		article.CategoryID = category.ID
		article.Category = *category
		article.Tags = tags
		article.WordCount = utf8.RuneCountInString(article.Content)
		logger.FromContext(c).Debug("更新文章业务", zap.Any("文章", article))
		return a.Repo.UpdateArticle(c, article)
	})
}

func (a *articleService) Create(c context.Context, article *model.Article) error {
	logger.FromContext(c).Debug("创建文章业务")
	return a.Ts.Transaction(c, func(c context.Context) error {
		// 1. 同步分类
		logger.FromContext(c).Debug("同步分类")
		category, err := a.CategoryRepo.SyncCategory(c, article.Category.Name)
		if err != nil {
			return err
		}
		// 2. 同步标签
		logger.FromContext(c).Debug("同步标签")
		tagsName := make([]string, len(article.Tags))
		for i, tag := range article.Tags {
			tagsName[i] = tag.Name
		}
		tags, err := a.TagRepo.SyncTags(c, tagsName)
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
			WordCount:  utf8.RuneCountInString(article.Content),
		}
		logger.FromContext(c).Debug("即将插入文章", zap.Any("文章", article))
		// 4. 调用 Create 方法
		return a.Repo.CreateArticle(c, article)
	})
}

func (a *articleService) Get(c context.Context, id int64) (model.Article, error) {
	return a.Repo.FindById(c, id, "Category", "Tags")
}

func (a *articleService) Delete(c context.Context, id int64) error {
	// 可以在这里增加删除后的逻辑（如清理 Redis）
	// return a.Repo.Delete(c, id)
	return a.Ts.Transaction(c, func(c context.Context) error {
		article, err := a.Repo.FindById(c, id, "Category", "Tags")
		if err != nil {
			return err
		}
		count, err := a.Repo.CountArticleByCategoryID(c, article.CategoryID)
		if err != nil {
			return err
		}
		if count <= 1 {
			logger.FromContext(c).Debug("删除文章业务 - 删除分类", zap.Int64("分类ID", article.CategoryID))
			if err := a.CategoryRepo.Delete(c, article.CategoryID); err != nil {
				return err
			}
		}
		tagIDs := make([]int64, len(article.Tags))
		for i, tag := range article.Tags {
			tagIDs[i] = tag.ID
		}
		counts, err := a.Repo.CountArticleByTagIDs(c, tagIDs)
		if err != nil {
			return err
		}
		// 删除标签（如果没有任何文章使用该标签）
		for tagID, count := range counts {
			if count <= 1 {
				logger.FromContext(c).Debug("删除文章业务 - 删除标签", zap.Int64("标签ID", tagID))
				if err := a.TagRepo.Delete(c, tagID); err != nil {
					return err
				}
			}
		}
		return a.Repo.Delete(c, id)
	})
}

func (a *articleService) List(c context.Context, page, pageSize int, filter *model.Article) ([]model.Article, int64, error) {
	return a.Repo.FindAllArticle(c, page, pageSize, filter)
}

func (a *articleService) Stats(c context.Context) (*model.ArticleStatsVO, error) {
	return a.Repo.GetArticleStats(c)
}
