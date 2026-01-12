package repository

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
	"gorm.io/gorm"
)

type BaseRepo[T any] interface {
	Conn(c context.Context) *gorm.DB
	Create(c context.Context, entity *T) error
	Delete(c context.Context, id uint) (int, error)
	FindAll(c context.Context, page int, pageSize int, entity *T, preloads ...string) ([]T, error)
	FindById(c context.Context, id uint, preloads ...string) (T, error)
}
type ArticleRepo interface {
	BaseRepo[model.Article]
	CreateArticle(c context.Context, article *model.Article) error
	FindAllArticle(c context.Context, page int, pageSize int, entity *model.Article) ([]model.Article, error)
	SyncCategory(ctx context.Context, name string) (*model.Category, error)
	SyncTags(ctx context.Context, names []string) ([]model.Tag, error)
	UpdateArticle(c context.Context, article *model.Article) error
}
type LoginRepo interface {
	CheckLogin(c context.Context, username string, password string) (uint, error)
}
