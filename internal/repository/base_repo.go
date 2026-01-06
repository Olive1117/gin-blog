package repository

import (
	"context"
	"errors"

	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BaseRepo[T any] struct {
	db *gorm.DB
}

func NewBaseRepo[T any](db *gorm.DB) *BaseRepo[T] {
	return &BaseRepo[T]{
		db: db,
	}
}

// GetById 根据id获取结构体，preloads是需要预加载的结构体字段
func (b *BaseRepo[T]) GetById(c context.Context, id uint, preloads ...string) (T, error) {
	db := b.Conn(c)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	var entity T
	err := db.Where("id = ?", id).First(&entity).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		logger.FromContext(c).Warn("记录未找到", zap.Uint("id", id))
		return entity, errs.ErrNotFound
	}
	return entity, nil
}

func (b *BaseRepo[T]) Create(c context.Context, entity *T) error {
	return gorm.G[T](b.Conn(c)).Create(c, entity)
}
func (b *BaseRepo[T]) Delete(c context.Context, id uint) (int, error) {
	return gorm.G[T](b.Conn(c)).Where("id = ?", id).Delete(c)
	// var entity T
	// return -1, b.Conn(c).Where("id = ?", id).Delete(&entity).Error
}
func (b *BaseRepo[T]) Update(c context.Context, id uint, data any) error {
	logger.FromContext(c).Debug("更新数据", zap.Uint("id", id), zap.Any("data", data))
	return b.Conn(c).Where("id = ?", id).Updates(data).Error
}
func (b *BaseRepo[T]) Conn(c context.Context) *gorm.DB {
	if tx, ok := c.Value("tx").(*gorm.DB); ok {
		return tx
	}
	return b.db.WithContext(c)
}
