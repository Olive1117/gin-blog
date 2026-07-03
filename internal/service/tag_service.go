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

func (ts *tagService) Create(c context.Context, tag *model.Tag) error {
	if id, err := ts.Repo.ExistByName(c, tag.Name); err != nil || id != 0 {
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
func (ts *tagService) List(c context.Context, que model.PageQuery, filter *model.Tag) (model.PageResult[model.Tag], error) {
	//TODO 这里应该写模糊查询，需要改baseRepo
	return ts.Repo.List(c, que, filter)
}
func (ts *tagService) Update(c context.Context, tag *model.Tag, id int64) error {
	existId, err := ts.Repo.ExistByName(c, tag.Name)
	if err != nil {
		return err
	}
	// 如果name存在，而且不是用户自己的，爆冲突警告
	if existId != 0 && existId != id {
		return errs.ErrExistTag
	}
	return ts.Repo.Update(c, id, tag)
}

func NewTagService(repo repository.TagRepo) TagService {
	return &tagService{
		Repo: repo,
	}
}
