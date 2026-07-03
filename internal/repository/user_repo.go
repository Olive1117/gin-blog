package repository

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
	"gorm.io/gorm"
)

type userRepo struct {
	BaseRepo[model.User]
}

// List implements [UserRepo].
func (u *userRepo) List(ctx context.Context, que model.PageQuery, filter *model.User) (model.PageResult[model.User], error) {
	return Paginate[model.User](u.Conn(ctx), que)
}

func (u *userRepo) FindByUniqueKeys(ctx context.Context, username string, email string) ([]model.User, error) {
	var users []model.User
	err := u.Conn(ctx).Select("id", "username", "email").Where("username = ? OR email = ?", username, email).Find(&users).Error
	return users, err
}
func (u *userRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := u.Conn(ctx).Model(&model.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		return &model.User{}, err
	}
	return &user, nil
}
func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{
		BaseRepo: NewBaseRepo[model.User](db),
	}
}
