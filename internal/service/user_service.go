package service

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/repository"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/spf13/cast"
)

type userService struct {
	Repo        repository.UserRepo
	ArticleRepo repository.ArticleRepo
}

func NewUserService(repo repository.UserRepo, articleRepo repository.ArticleRepo) UserService {
	return &userService{
		Repo:        repo,
		ArticleRepo: articleRepo,
	}
}

func (ts *userService) Create(c context.Context, user *model.User) error {
	if users, err := ts.Repo.FindByUniqueKeys(c, user.Username, user.Email); err != nil || len(users) > 0 {
		// 唯一性信息冲突
		for _, u := range users {
			if u.Username == user.Username {
				return errs.ErrExistUsername
			}
			if u.Email == user.Email {
				return errs.ErrExistEmail
			}
		}
		return err
	}
	return ts.Repo.Create(c, user)
}
func (ts *userService) Delete(c context.Context, id int64) error {
	return ts.Repo.Delete(c, id)
}
func (ts *userService) Get(c context.Context, id int64) (model.User, error) {
	user, err := ts.Repo.FindById(c, id)
	if err != nil {
		return model.User{}, err
	}
	postCount, err := ts.ArticleRepo.CountArticleByUserID(c, id)
	if err != nil {
		return model.User{}, err
	}
	user.PostCount = cast.ToInt(postCount)
	return user, nil
}
func (ts *userService) List(c context.Context, page int, pageSize int, filter *model.User) ([]model.User, int64, error) {
	return ts.Repo.FindAll(c, page, pageSize, filter)
}
func (ts *userService) Update(c context.Context, user *model.User, id int64) error {
	users, err := ts.Repo.FindByUniqueKeys(c, user.Username, user.Email)
	if err != nil {
		return err
	}
	for _, u := range users {
		// 如果不是用户自己的信息，发生冲突
		if u.ID != id {
			if u.Username == user.Username {
				return errs.ErrExistUsername
			}
			if u.Email == user.Email {
				return errs.ErrExistEmail
			}
		}
	}
	return ts.Repo.Update(c, id, user)
}
