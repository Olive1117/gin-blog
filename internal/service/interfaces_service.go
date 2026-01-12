package service

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
)

type ArticleService interface {
	Create(c context.Context, article *model.Article) error
	Delete(c context.Context, id uint) (int, error)
	Get(c context.Context, id uint) (model.Article, error)
	List(c context.Context, page int, pageSize int, filter *model.Article) ([]model.Article, error)
	Update(c context.Context, article *model.Article) error
}
type LoginService interface {
	Login(c context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
}
