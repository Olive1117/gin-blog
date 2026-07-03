package repository

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
	"gorm.io/gorm"
)

type friendlinkRepo struct {
	BaseRepo[model.FriendLink]
}

// List implements [FriendLinkRepo].
func (f *friendlinkRepo) List(ctx context.Context, que model.PageQuery) (model.PageResult[model.FriendLink], error) {
	return Paginate[model.FriendLink](f.Conn(ctx), que)
}

func NewFriendLinkRepo(db *gorm.DB) FriendLinkRepo {
	return &friendlinkRepo{
		BaseRepo: NewBaseRepo[model.FriendLink](db),
	}
}
