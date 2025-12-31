package repository

import (
	"context"
	"errors"

	"github.com/Olive1117/gin-blog/pkg/errs"
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
func (b *BaseRepo[T]) GetById(c context.Context, id uint, preloads ...string) (T, error) {
	db := b.Conn(c)
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	var entity T
	err := db.Where("id = ?", id).First(&entity).Error
	// entity, err := gorm.G[T](b.Conn(c)).Where("id = ?", id).First(c)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return entity, errs.ErrNotFound
	}
	return entity, nil
}

func (b *BaseRepo[T]) Create(c context.Context, entity *T) error {
	return gorm.G[T](b.Conn(c)).Create(c, entity)
}
func (b *BaseRepo[T]) Delete(c context.Context, id uint) (int, error) {
	return gorm.G[T](b.Conn(c)).Where("id = ?", id).Delete(c)
}
func (b *BaseRepo[T]) Update(c context.Context, id uint, data T) (int, error) {
	return gorm.G[T](b.Conn(c)).Where("id = ?", id).Updates(c, data)
}
func (b *BaseRepo[T]) Conn(c context.Context) *gorm.DB {
	if tx, ok := c.Value("tx").(*gorm.DB); ok {
		return tx
	}
	return b.db.WithContext(c)
}
