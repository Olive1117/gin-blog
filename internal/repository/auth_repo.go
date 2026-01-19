package repository

import (
	"context"
	"errors"

	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"gorm.io/gorm"
)

type authRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) AuthRepo {
	return &authRepo{
		db: db,
	}
}

func (l *authRepo) CheckAuth(c context.Context, username string, password string) (int64, error) {
	var auth model.User
	err := l.db.WithContext(c).Select("id").Where(&model.User{Username: username, Password: password}).First(&auth).Error
	if auth.ID > 0 {
		return auth.ID, nil
	} else {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errs.ErrAuth
		}
		return 0, err
	}
}
