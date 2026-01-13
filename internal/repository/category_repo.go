package repository

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
	"gorm.io/gorm"
)

type categoryRepo struct {
	BaseRepo[model.Category]
}

func NewCategoryRepo(db *gorm.DB) CategoryRepo {
	return &categoryRepo{
		BaseRepo: NewBaseRepo[model.Category](db),
	}
}

// 同步分类
func (r *categoryRepo) SyncCategory(ctx context.Context, name string) (*model.Category, error) {
	var category = &model.Category{Name: name}
	err := r.Conn(ctx).WithContext(ctx).Where("name = ?", name).FirstOrCreate(category).Error
	return category, err
}
