package repository

import (
	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/service"
	"gorm.io/gorm"
)

var _ service.LoginStore = (*LoginRepo)(nil)

type LoginRepo struct {
	db *gorm.DB
}

func NewLoginRepo(db *gorm.DB) *LoginRepo {
	return &LoginRepo{
		db: db,
	}
}

func (l *LoginRepo) CheckLogin(username string, password string) (uint, error) {
	var login model.Login
	err := l.db.Select("id").Where(&model.Login{Username: username, Password: password}).First(&login).Error
	if login.ID > 0 {
		return login.ID, err
	} else {
		return 0, err
	}
}
