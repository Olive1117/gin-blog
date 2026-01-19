package service

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/repository"
	"github.com/Olive1117/gin-blog/pkg/errs"
)

type categoryService struct {
	Repo repository.CategoryRepo
}

func NewCategoryService(repo repository.CategoryRepo) CategoryService {
	return &categoryService{
		Repo: repo,
	}
}

func (cs *categoryService) Create(ctx context.Context, category *model.Category) error {
	if ok, err := cs.Repo.ExistByName(ctx, category.Name); err != nil || ok {
		return errs.ErrExistCategory
	}
	return cs.Repo.Create(ctx, category)
}
func (cs *categoryService) Delete(ctx context.Context, id int64) error {
	return cs.Repo.Delete(ctx, id)
}
func (cs *categoryService) Get(ctx context.Context, id int64) (model.Category, error) {
	return cs.Repo.FindById(ctx, id)
}
func (cs *categoryService) List(ctx context.Context, page int, pageSize int, filter *model.Category) ([]model.Category, int64, error) {
	//TODO 这里应该写模糊查询，需要改baseRepo
	return cs.Repo.FindAll(ctx, page, pageSize, filter)
}
func (cs *categoryService) Update(ctx context.Context, category *model.Category) error {
	if ok, err := cs.Repo.ExistByName(ctx, category.Name); err != nil || ok {
		return errs.ErrExistCategory
	}
	return cs.Repo.Update(ctx, category.ID, category)
}
