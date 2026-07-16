package model

import (
	"time"

	"github.com/Olive1117/gin-blog/pkg/idgen"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int64          `gorm:"primaryKey;autoIncrement:false" json:"id,string"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index:,composite:deletedat" json:"-"`
	CreatedBy int64          `gorm:"default:0;comment:创建者ID"`
	UpdatedBy int64          `gorm:"default:0;comment:修改者ID"`
	DeletedBy int64          `gorm:"default:0;comment:删除者ID"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == 0 {
		b.ID = idgen.NextID()
	}
	return
}

type PageQuery struct {
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}
type PageResult[T any] struct {
	List  []T
	Total int64
}

func (q PageQuery) Offset() int {
	if q.Page < 1 {
		return 0
	}
	return (q.Page - 1) * q.PageSize
}

func (q PageQuery) Limit() int {
	if q.PageSize < 1 {
		return 10 // 默认
	}
	if q.PageSize > 100 {
		return 100 // 最大
	}
	return q.PageSize
}

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type PageResponse[T any] struct {
	List     []T   `json:"list"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

type RefreshCookie struct {
	RefreshToken string
	ExpiresAt    time.Time
}
