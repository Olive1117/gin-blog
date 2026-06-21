package service

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
)

type BaseService[T any] interface {
	Create(c context.Context, entity *T) error
	Delete(c context.Context, id int64) error
	Get(c context.Context, id int64) (T, error)
	List(c context.Context, page int, pageSize int, filter *T) ([]T, int64, error)
	Update(c context.Context, entity *T, id int64) error
}
type ArticleService interface {
	BaseService[model.Article]
	Stats(c context.Context) (*model.ArticleStatsDTO, error)
}
type AuthService interface {
	Auth(c context.Context, req *model.AuthRequest) (*model.AuthResponse, error)
}
type CategoryService interface {
	BaseService[model.Category]
}
type TagService interface {
	BaseService[model.Tag]
}
type UserService interface {
	BaseService[model.User]
}
