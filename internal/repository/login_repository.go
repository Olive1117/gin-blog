package repository

import (
	"github.com/Olive1117/gin-blog/model"
	"gorm.io/gorm"
)

type loginRepo struct {
	db *gorm.DB
}

func NewLoginRepo(db *gorm.DB) *loginRepo {
	return &loginRepo{
		db: db,
	}
}

func (l *loginRepo) CheckLogin(username string, password string) (uint, error) {
	var login model.Login
	err := l.db.Select("id").Where(&model.Login{Username: username, Password: password}).First(&login).Error
	if login.ID > 0 {
		return login.ID, err
	} else {
		return 0, err
	}
}
