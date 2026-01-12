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

// 同步分类
func (r *articleRepo) SyncCategory(ctx context.Context, name string) (*model.Category, error) {
	var category = &model.Category{Name: name}
	err := r.Conn(ctx).WithContext(ctx).Where("name = ?", name).FirstOrCreate(category).Error
	return category, err
}

// 同步标签
func (r *articleRepo) SyncTags(ctx context.Context, names []string) ([]model.Tag, error) {
	if len(names) == 0 {
		return []model.Tag{}, nil
	}
	var tags []model.Tag
	// 根据标签名列表获取已有标签
	if err := r.Conn(ctx).WithContext(ctx).Where("name IN ?", names).Find(&tags).Error; err != nil {
		return nil, err
	}
	// 找出不存在的标签并创建
	existingNames := make(map[string]bool)
	for _, tag := range tags {
		existingNames[tag.Name] = true
	}
	var newTags []model.Tag
	for _, name := range names {
		if !existingNames[name] {
			newTags = append(newTags, model.Tag{Name: name})
		}
	}
	if len(newTags) > 0 {
		if err := r.Conn(ctx).WithContext(ctx).Create(&newTags).Error; err != nil {
			return nil, err
		}
		tags = append(tags, newTags...)
	}
	return tags, nil
}

func (r *articleRepo) CreateArticle(c context.Context, article *model.Article) error {
	// 业务代码中已经处理了标签的关联关系，这里创建文章时忽略标签字段
	return r.Conn(c).Omit("Tags.*").Create(article).Error
}

func (r *articleRepo) FindAllArticle(c context.Context, page, pageSize int, entity *model.Article) ([]model.Article, error) {
	db := r.Conn(c)
	if entity.Category.Name != "" {
		db = db.Joins("Category").Where("category.name = ?", entity.Category.Name)
	}
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
	if entity != nil {
		db = db.Where(entity)
	}
	var articles []model.Article
	offset := (page - 1) * pageSize
	err := db.Preload("Category").Preload("Tags").
		Offset(offset).Limit(pageSize).
		Find(&articles).Error
	return articles, err
}

func (r *articleRepo) UpdateArticle(c context.Context, article *model.Article) error {
	return r.Conn(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Tags", "Category").Where("id = ?", article.ID).Updates(article).Error; err != nil {
			return err
		}
		return tx.Model(article).Omit("Tags.*").Association("Tags").Replace(article.Tags)
	})
}
