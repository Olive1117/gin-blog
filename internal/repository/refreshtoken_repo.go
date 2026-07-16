package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"gorm.io/gorm"
)

type refreshtokenRepo struct {
	DBConn
}

// CleanupExpired implements [RefreshTokenRepo].
func (r *refreshtokenRepo) CleanupExpired(c context.Context, before time.Time) error {
	return r.Conn(c).Where("revoked_at is not null").Where("expires_at < ?", before).Delete(&model.RefreshTokens{}).Error
}

// Create implements [RefreshTokenRepo].
func (r *refreshtokenRepo) Create(c context.Context, value model.RefreshTokens) error {
	return r.Conn(c).Create(&value).Error
}

// DeleteByJTI implements [RefreshTokenRepo].
func (r *refreshtokenRepo) DeleteByJTI(c context.Context, jti string) error {
	return r.Conn(c).Delete(&model.RefreshTokens{Jti: jti}).Error
}

// GetByJTI implements [RefreshTokenRepo].
func (r *refreshtokenRepo) GetByJTI(c context.Context, jti string) (model.RefreshTokens, error) {
	var res model.RefreshTokens
	err := r.Conn(c).Where("jti = ?", jti).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.WarnContext(c, "记录未找到", logger.String("jti", jti))
			return res, errs.ErrNotFound
		}
		logger.ErrorContext(c, "获取记录失败", logger.Err(err))
		return res, err
	}
	return res, nil
}

// RevokeByJTI implements [RefreshTokenRepo].
func (r *refreshtokenRepo) RevokeByJTI(c context.Context, jti string) error {
	return r.Conn(c).Model(&model.RefreshTokens{}).Where("jti = ?", jti).Update("revoked_at", time.Now()).Error
}

// RevokeByUser implements [RefreshTokenRepo].
func (r *refreshtokenRepo) RevokeByUser(c context.Context, userid int64) error {
	return r.Conn(c).Model(&model.RefreshTokens{}).Where("user_id = ?", userid).Update("revoked_at", time.Now()).Error
}

func NewreFreshTokenRepo(db *gorm.DB) RefreshTokenRepo {
	return &refreshtokenRepo{
		DBConn: NewBaseRepo[model.RefreshTokens](db),
	}
}
