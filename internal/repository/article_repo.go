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

// 同步分类
func (r *ArticleRepo) SyncCategory(ctx context.Context, name string) (*model.Category, error) {
	var category = &model.Category{Name: name}
	err := r.Conn(ctx).WithContext(ctx).Where("name = ?", name).FirstOrCreate(category).Error
	return category, err
}

// 同步标签
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

func (r *ArticleRepo) FindAllArticle(c context.Context, page, pageSize int, entity *model.Article) ([]model.Article, error) {
	db := r.Conn(c)
	if entity.Category != (model.Category{}) {
		db = db.Joins("Category").Where("categories.name = ?", entity.Category.Name)
	}
	if len(entity.Tags) > 0 {
		tagNames := make([]string, len(entity.Tags))
		for i, tag := range entity.Tags {
			tagNames[i] = tag.Name
		}
		db = db.Joins("JOIN blog_article_tag ON blog_article_tag.article_id = blog_article.id").
			Joins("JOIN blog_tag ON blog_tag.id = blog_article_tag.tag_id").
			Where("blog_tag.name IN ?", tagNames)
	}
	if entity != nil {
		db = db.Where(entity)
	}
	var articles []model.Article
	offset := (page - 1) * pageSize
	err := db.Preload("Category").Preload("Tags").Offset(offset).Limit(pageSize).Find(&articles).Error
	return articles, err
}
