package repository

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
	"gorm.io/gorm"
)

type ArticleRepo struct {
	*BaseRepo[model.Article]
}

func NewArticleRepo(db *gorm.DB) *ArticleRepo {
	return &ArticleRepo{
		BaseRepo: NewBaseRepo[model.Article](db),
	}
}

func (r *ArticleRepo) SyncCategory(ctx context.Context, name string) (*model.Category, error) {
	var category = &model.Category{Name: name}
	err := r.Conn(ctx).WithContext(ctx).Where("name = ?", name).FirstOrCreate(category).Error
	return category, err
}

func (r *ArticleRepo) SyncTags(ctx context.Context, names []string) ([]model.Tag, error) {
	var tags []model.Tag
	for _, name := range names {
		var tag = &model.Tag{Name: name}
		if err := r.Conn(ctx).WithContext(ctx).Where("name = ?", name).FirstOrCreate(&tag).Error; err != nil {
			return nil, err
		}
		tags = append(tags, *tag)
	}
	return tags, nil
}

func (r *ArticleRepo) CreateArticle(c context.Context, article *model.Article) error {
	return r.Conn(c).Omit("Tags.*").Create(article).Error
}

func (r *ArticleRepo) FindByID(c context.Context, id uint) (model.Article, error) {
	// var article model.Article
	// err := r.Conn(c).Preload("Category").Preload("Tags").Where("id = ?", id).First(&article).Error
	return gorm.G[model.Article](r.Conn(c)).Preload("Category", nil).Preload("Tags", nil).Where("id = ?", id).First(c)
}
