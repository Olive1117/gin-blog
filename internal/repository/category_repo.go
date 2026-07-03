package repository

import (
	"context"
	"errors"

	"github.com/Olive1117/gin-blog/internal/model"
	"gorm.io/gorm"
)

type categoryRepo struct {
	BaseRepo[model.Category]
}

// List implements [CategoryRepo].
func (r *categoryRepo) List(ctx context.Context, que model.PageQuery, filter *model.Category) (model.PageResult[model.Category], error) {
	db := r.Conn(ctx).Where(filter)
	return Paginate[model.Category](db, que)
}

// 同步分类
func (r *categoryRepo) SyncCategory(ctx context.Context, name string) (*model.Category, error) {
	var category = &model.Category{Name: name}
	err := r.Conn(ctx).Where("name = ?", name).FirstOrCreate(category).Error
	return category, err
}

func (r *categoryRepo) ExistByName(ctx context.Context, name string) (int64, error) {
	var category model.Category
	err := r.Conn(ctx).Select("id").Where("name = ?", name).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil // 表示不存在
		}
		return 0, err // 表示数据库报错
	}
	return category.ID, nil
}

func NewCategoryRepo(db *gorm.DB) CategoryRepo {
	return &categoryRepo{
		BaseRepo: NewBaseRepo[model.Category](db),
	}
}
