package repository

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
	"gorm.io/gorm"
)

type articleRepo struct {
	BaseRepo[model.Article]
}

func NewArticleRepo(db *gorm.DB) ArticleRepo {
	return &articleRepo{
		BaseRepo: NewBaseRepo[model.Article](db),
	}
}

func (r *articleRepo) CreateArticle(c context.Context, article *model.Article) error {
	// 业务代码中已经处理了标签的关联关系，这里创建文章时忽略标签字段
	return r.Conn(c).Omit("Tags.*").Create(article).Error
}

func (r *articleRepo) FindAllArticle(c context.Context, page, pageSize int, entity *model.Article) ([]model.Article, int64, error) {
	db := r.Conn(c)
	if entity.Category.Name != "" {
		db = db.Joins("Category").Where("category.name = ?", entity.Category.Name)
	}
	// 构建查询条件
	if len(entity.Tags) > 0 {
		tagNames := make([]string, len(entity.Tags))
		for i, tag := range entity.Tags {
			tagNames[i] = tag.Name
		}
		// 通过子查询过滤包含指定标签的文章
		subQuery := r.Conn(c).Table("blog_article_tag").
			Select("article_id").
			Joins("JOIN blog_tag ON blog_tag.id = blog_article_tag.tag_id").
			Where("blog_tag.name IN ?", tagNames)
		db = db.Where("blog_article.id IN (?)", subQuery)
	}
	// 基础过滤条件
	if entity != nil {
		db = db.Where(entity)
	}
	var tatol int64
	if err := db.Model(&model.Article{}).Count(&tatol).Error; err != nil {
		return nil, 0, err
	}
	var articles []model.Article
	offset := (page - 1) * pageSize
	if err := db.Preload("Category").Preload("Tags").
		Omit("Content").
		Offset(offset).Limit(pageSize).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}
	return articles, tatol, nil
}
func (r *articleRepo) UpdateArticle(c context.Context, article *model.Article) error {
	return r.Conn(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Tags", "Category").Where("id = ?", article.ID).Updates(article).Error; err != nil {
			return err
		}
		return tx.Model(article).Omit("Tags.*").Association("Tags").Replace(article.Tags)
	})
}
