package service

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/repository"
	"github.com/Olive1117/gin-blog/pkg/errs"
)

type tagService struct {
	Repo repository.TagRepo
}

func NewTagService(repo repository.TagRepo) TagService {
	return &tagService{
		Repo: repo,
	}
}

func (ts *tagService) Create(c context.Context, tag *model.Tag) error {
	if ok, err := ts.Repo.ExistByName(c, tag.Name); err != nil || ok {
		return errs.ErrExistTag
	}
	return ts.Repo.Create(c, tag)
}
func (ts *tagService) Delete(c context.Context, id int64) error {
	return ts.Repo.Delete(c, id)
}
func (ts *tagService) Get(c context.Context, id int64) (model.Tag, error) {
	return ts.Repo.FindById(c, id)
}
func (ts *tagService) List(c context.Context, page int, pageSize int, filter *model.Tag) ([]model.Tag, int64, error) {
	//TODO 这里应该写模糊查询，需要改baseRepo
	return ts.Repo.FindAll(c, page, pageSize, filter)
}
func (ts *tagService) Update(c context.Context, tag *model.Tag) error {
	if ok, err := ts.Repo.ExistByName(c, tag.Name); err != nil || ok {
		return errs.ErrExistTag
	}
	return ts.Repo.Update(c, tag.ID, tag)
}
