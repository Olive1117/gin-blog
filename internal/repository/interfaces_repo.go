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
	FindById(c context.Context, id int64, preloads ...string) (T, error)
	Update(c context.Context, id int64, data *T) error
}
type ArticleRepo interface {
	BaseRepo[model.Article]
	CreateArticle(c context.Context, article *model.Article) error
	List(c context.Context, que model.PageQuery, entity *model.Article) (model.PageResult[model.Article], error)
	UpdateArticle(c context.Context, article *model.Article) error
	CountArticleByUserID(c context.Context, userID int64) (int64, error)
	CountArticleByCategoryID(c context.Context, categoryID int64) (int64, error)
	CountArticleByTagIDs(c context.Context, tagIDs []int64) (map[int64]int64, error)
	GetArticleStats(c context.Context) (*model.ArticleStatsVO, error)
}
type TagRepo interface {
	BaseRepo[model.Tag]
	SyncTags(ctx context.Context, names []string) ([]model.Tag, error)
	ExistByName(ctx context.Context, name string) (int64, error)
	List(ctx context.Context, que model.PageQuery, filter *model.Tag) (model.PageResult[model.Tag], error)
}
type CategoryRepo interface {
	BaseRepo[model.Category]
	SyncCategory(ctx context.Context, name string) (*model.Category, error)
	ExistByName(ctx context.Context, name string) (int64, error)
	List(ctx context.Context, que model.PageQuery, filter *model.Category) (model.PageResult[model.Category], error)
}
type UserRepo interface {
	BaseRepo[model.User]
	FindByUniqueKeys(ctx context.Context, username string, email string) ([]model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	List(ctx context.Context, que model.PageQuery, filter *model.User) (model.PageResult[model.User], error)
}
type FriendLinkRepo interface {
	BaseRepo[model.FriendLink]
	List(context.Context, model.PageQuery) (model.PageResult[model.FriendLink], error)
}
