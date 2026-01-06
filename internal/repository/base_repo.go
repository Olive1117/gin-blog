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

// FindAll 根据条件获取多个记录，支持分页和预加载
func (b *BaseRepo[T]) FindAll(c context.Context, page, pageSize int, entity *T, preloads ...string) ([]T, error) {
	db := b.Conn(c)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	if entity != nil {
		db = db.Where(entity)
	}
	var entities []T
	offset := (page - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Find(&entities).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.FromContext(c).Warn("记录未找到")
			return entities, errs.ErrNotFound
		}
		logger.FromContext(c).Error("获取记录失败", zap.Error(err))
		return entities, err
	}
	return entities, nil
}

// FindById 根据id获取结构体，preloads是需要预加载的结构体字段
func (b *BaseRepo[T]) FindById(c context.Context, id uint, preloads ...string) (T, error) {
	db := b.Conn(c)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	var entity T
	err := db.Where("id = ?", id).First(&entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.FromContext(c).Warn("记录未找到", zap.Uint("id", id))
			return entity, errs.ErrNotFound
		}
		logger.FromContext(c).Error("获取记录失败", zap.Error(err))
		return entity, err
	}
	return entity, nil
}

func (b *BaseRepo[T]) Create(c context.Context, entity *T) error {
	return gorm.G[T](b.Conn(c)).Create(c, entity)
}

func (b *BaseRepo[T]) Delete(c context.Context, id uint) (int, error) {
	return gorm.G[T](b.Conn(c)).Where("id = ?", id).Delete(c)
}

func (b *BaseRepo[T]) Update(c context.Context, id uint, data any) error {
	return b.Conn(c).Where("id = ?", id).Updates(data).Error
}

func (b *BaseRepo[T]) Conn(c context.Context) *gorm.DB {
	if tx, ok := c.Value("tx").(*gorm.DB); ok {
		return tx
	}
	return b.db.WithContext(c)
}
