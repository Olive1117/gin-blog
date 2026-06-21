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
		Order("created_at DESC").
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

func (r *articleRepo) CountArticleByUserID(c context.Context, userID int64) (int64, error) {
	var count int64
	if err := r.Conn(c).Model(&model.Article{}).Where("created_by = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *articleRepo) CountArticleByCategoryID(c context.Context, categoryID int64) (int64, error) {
	var count int64
	if err := r.Conn(c).Model(&model.Article{}).Where("category_id = ?", categoryID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *articleRepo) CountArticleByTagIDs(c context.Context, tagIDs []int64) (map[int64]int64, error) {
	type result struct {
		TagID int64
		Count int64
	}
	var rows []result
	if err := r.Conn(c).Raw("SELECT tag_id, COUNT(*) as count FROM blog_article_tag WHERE tag_id IN ? GROUP BY tag_id", tagIDs).Scan(&rows).Error; err != nil {
		return nil, err
	}
	counts := make(map[int64]int64, len(rows))
	for _, r := range rows {
		counts[r.TagID] = r.Count
	}
	return counts, nil
}

func (r *articleRepo) GetArticleStats(c context.Context) (*model.ArticleStatsDTO, error) {
	var stats model.ArticleStatsDTO
	stats.TotalByCategory = make(map[string]int64)
	stats.TotalByTag = make(map[string]int64)
	if err := r.Conn(c).Model(&model.Article{}).Where("blog_article.state = ?", 1).Count(&stats.Total).Error; err != nil {
		return nil, err
	}
	var categoryRows []struct {
		Name  string
		Count int64
	}
	if err := r.Conn(c).Raw(`
			SELECT c.name, COUNT(*) as count
			FROM blog_article a
			JOIN blog_category c ON a.category_id = c.id
			WHERE a.state = 1 AND a.deleted_at IS NULL
			GROUP BY c.id, c.name
		`).
		Scan(&categoryRows).Error; err != nil {
		return nil, err
	}
	for _, row := range categoryRows {
		stats.TotalByCategory[row.Name] = row.Count
	}
	var tagRows []struct {
		Name  string
		Count int64
	}
	if err := r.Conn(c).Raw(`
		SELECT t.name, COUNT(*) as count
		FROM blog_article_tag at
		JOIN blog_tag t ON at.tag_id = t.id
		JOIN blog_article a ON at.article_id = a.id
		WHERE a.state = 1 AND a.deleted_at IS NULL
		GROUP BY t.name
	`).Scan(&tagRows).Error; err != nil {
		return nil, err
	}
	for _, row := range tagRows {
		stats.TotalByTag[row.Name] = row.Count
	}

	return &stats, nil
}
