package repository

import (
	"context"
	"errors"

	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type baseRepo[T any] struct {
	db *gorm.DB
}

func NewBaseRepo[T any](db *gorm.DB) BaseRepo[T] {
	return &baseRepo[T]{
		db: db,
	}
}

// FindAll 根据条件获取多个记录，支持分页和预加载
func (b *baseRepo[T]) FindAll(c context.Context, page, pageSize int, entity *T, preloads ...string) ([]T, int64, error) {
	db := b.Conn(c)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	if entity != nil {
		db = db.Where(entity)
	}
	var total int64
	if err := db.Model(new(T)).Count(&total).Error; err != nil {
		logger.FromContext(c).Error("获取记录总数失败", zap.Error(err))
		return nil, 0, err
	}
	var entities []T
	offset := (page - 1) * pageSize
	err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.FromContext(c).Warn("记录未找到")
			return entities, total, errs.ErrNotFound
		}
		logger.FromContext(c).Error("获取记录失败", zap.Error(err))
		return nil, 0, err
	}
	return entities, total, nil
}

// FindById 根据id获取结构体，preloads是需要预加载的结构体字段
func (b *baseRepo[T]) FindById(c context.Context, id int64, preloads ...string) (T, error) {
	db := b.Conn(c)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	var entity T
	err := db.Where("id = ?", id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.FromContext(c).Warn("记录未找到", zap.Int64("id", id))
			return entity, errs.ErrNotFound
		}
		logger.FromContext(c).Error("获取记录失败", zap.Error(err))
		return entity, err
	}
	return entity, nil
}

// Create 创建结构体数据，复杂结构体请使用具体的Repo方法
func (b *baseRepo[T]) Create(c context.Context, entity *T) error {
	return gorm.G[T](b.Conn(c)).Create(c, entity)
}

// Delete 删除结构体数据，返回删除的行数
func (b *baseRepo[T]) Delete(c context.Context, id int64) error {
	var model T
	return b.Conn(c).Transaction(func(tx *gorm.DB) error {
		return tx.Where("id = ?", id).Delete(&model).Error
	})
}

// Update 更新结构体数据，复杂结构体请使用具体的Repo方法
func (b *baseRepo[T]) Update(c context.Context, id int64, data any) error {
	return b.Conn(c).Where("id = ?", id).Updates(data).Error
}

// Conn 事务获取
func (b *baseRepo[T]) Conn(c context.Context) *gorm.DB {
	if tx, ok := c.Value("tx").(*gorm.DB); ok {
		return tx
	}
	return b.db.WithContext(c)
}
