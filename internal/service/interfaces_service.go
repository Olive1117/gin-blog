package service

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
)

type BaseService[T any] interface {
	Create(c context.Context, entity *T) error
	Delete(c context.Context, id int64) error
	Get(c context.Context, id int64) (T, error)
	Update(c context.Context, entity *T, id int64) error
}
type ArticleService interface {
	BaseService[model.Article]
	Stats(c context.Context) (*model.ArticleStatsVO, error)
	List(c context.Context, que model.PageQuery, filter *model.Article) (model.PageResult[model.Article], error)
}
type CategoryService interface {
	BaseService[model.Category]
	List(ctx context.Context, que model.PageQuery, filter *model.Category) (model.PageResult[model.Category], error)
}
type TagService interface {
	BaseService[model.Tag]
	List(ctx context.Context, que model.PageQuery, filter *model.Tag) (model.PageResult[model.Tag], error)
}
type UserService interface {
	BaseService[model.User]
	Login(c context.Context, req model.LoginRequest) (model.AuthResponse, error)
	ChangePassword(c context.Context, username string, oldPassword string, newPassword string) error
	List(ctx context.Context, que model.PageQuery, filter *model.User) (model.PageResult[model.User], error)
}
type FriendLinkService interface {
	BaseService[model.FriendLink]
	List(context.Context, model.PageQuery) (model.PageResult[model.FriendLink], error)
}
