package service

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/repository"
)

type friendlinkService struct {
	repo repository.FriendLinkRepo
}

// Create implements [FriendLinkService].
func (f *friendlinkService) Create(c context.Context, entity *model.FriendLink) error {
	return f.repo.Create(c, entity)
}

// Delete implements [FriendLinkService].
func (f *friendlinkService) Delete(c context.Context, id int64) error {
	return f.repo.Delete(c, id)
}

// Get implements [FriendLinkService].
func (f *friendlinkService) Get(c context.Context, id int64) (model.FriendLink, error) {
	return f.repo.FindById(c, id)
}

// List implements [FriendLinkService].
func (f *friendlinkService) List(c context.Context, que model.PageQuery) (model.PageResult[model.FriendLink], error) {
	return f.repo.List(c, que)
}

// Update implements [FriendLinkService].
func (f *friendlinkService) Update(c context.Context, entity *model.FriendLink, id int64) error {
	return f.repo.Update(c, id, entity)
}

func NewFriendLinkService(repo repository.FriendLinkRepo) FriendLinkService {
	return &friendlinkService{repo: repo}
}
