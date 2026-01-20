package repository

import (
	"context"

	"github.com/Olive1117/gin-blog/internal/model"
	"gorm.io/gorm"
)

type userRepo struct {
	BaseRepo[model.User]
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{
		BaseRepo: NewBaseRepo[model.User](db),
	}
}

func (u *userRepo) FindByUniqueKeys(ctx context.Context, username string, email string) ([]model.User, error) {
	var users []model.User
	err := u.Conn(ctx).Select("id", "username", "email").Where("username = ? OR email = ?", username, email).Find(&users).Error
	return users, err
}
