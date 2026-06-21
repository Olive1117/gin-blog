package repository

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
	"gorm.io/gorm"
)

type BaseRepo[T any] interface {
	Conn(c context.Context) *gorm.DB
	Create(c context.Context, entity *T) error
	Delete(c context.Context, id int64) error
	FindAll(c context.Context, page int, pageSize int, entity *T, preloads ...string) ([]T, int64, error)
	FindById(c context.Context, id int64, preloads ...string) (T, error)
	Update(c context.Context, id int64, data any) error
}
type ArticleRepo interface {
	BaseRepo[model.Article]
	CreateArticle(c context.Context, article *model.Article) error
	FindAllArticle(c context.Context, page int, pageSize int, entity *model.Article) ([]model.Article, int64, error)
	UpdateArticle(c context.Context, article *model.Article) error
	CountArticleByUserID(c context.Context, userID int64) (int64, error)
	CountArticleByCategoryID(c context.Context, categoryID int64) (int64, error)
	CountArticleByTagIDs(c context.Context, tagIDs []int64) (map[int64]int64, error)
	GetArticleStats(c context.Context) (*model.ArticleStatsDTO, error)
}
type AuthRepo interface {
	CheckAuth(c context.Context, username string, password string) (int64, error)
}
type TagRepo interface {
	BaseRepo[model.Tag]
	SyncTags(ctx context.Context, names []string) ([]model.Tag, error)
	ExistByName(ctx context.Context, name string) (int64, error)
}
type CategoryRepo interface {
	BaseRepo[model.Category]
	SyncCategory(ctx context.Context, name string) (*model.Category, error)
	ExistByName(ctx context.Context, name string) (int64, error)
}
type UserRepo interface {
	BaseRepo[model.User]
	FindByUniqueKeys(ctx context.Context, username string, email string) ([]model.User, error)
}
