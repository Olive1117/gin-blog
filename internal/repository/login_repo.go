package repository

import (
	"context"
	"errors"

	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"gorm.io/gorm"
)

type loginRepo struct {
	db *gorm.DB
}

func NewLoginRepo(db *gorm.DB) LoginRepo {
	return &loginRepo{
		db: db,
	}
}

func (l *loginRepo) CheckLogin(c context.Context, username string, password string) (uint, error) {
	var login model.User
	err := l.db.WithContext(c).Select("id").Where(&model.User{Username: username, Password: password}).First(&login).Error
	if login.ID > 0 {
		return login.ID, nil
	} else {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errs.ErrLogin
		}
		return 0, err
	}
}
