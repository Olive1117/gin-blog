package repository

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
	"gorm.io/gorm"
)

type tagRepo struct {
	BaseRepo[model.Tag]
}

func NewTagRepo(db *gorm.DB) TagRepo {
	return &tagRepo{
		BaseRepo: NewBaseRepo[model.Tag](db),
	}
}

// 同步标签
func (r *tagRepo) SyncTags(ctx context.Context, names []string) ([]model.Tag, error) {
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
