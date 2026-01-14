package service

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
)

type ArticleService interface {
	Create(c context.Context, article *model.Article) error
	Delete(c context.Context, id int64) (int, error)
	Get(c context.Context, id int64) (model.Article, error)
	List(c context.Context, page int, pageSize int, filter *model.Article) ([]model.Article, int64, error)
	Update(c context.Context, article *model.Article) error
}
type LoginService interface {
	Login(c context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
}
type CategoryService interface {
	Create(c context.Context, category *model.Category) error
	Delete(c context.Context, id int64) (int, error)
	Get(c context.Context, id int64) (model.Category, error)
	List(c context.Context, page int, pageSize int, filter *model.Category) ([]model.Category, int64, error)
	Update(c context.Context, category *model.Category) error
}
type TagService interface {
	Create(c context.Context, tag *model.Tag) error
	Delete(c context.Context, id int64) (int, error)
	Get(c context.Context, id int64) (model.Tag, error)
	List(c context.Context, page int, pageSize int, filter *model.Tag) ([]model.Tag, int64, error)
	Update(c context.Context, tag *model.Tag) error
}
