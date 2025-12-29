package model

import "gorm.io/gorm"

type BaseModel struct {
	gorm.Model

	CreatedBy uint64 `gorm:"column:created_by"`
	UpdatedBy uint64 `gorm:"column:updated_by"`
	DeletedBy uint64 `gorm:"column:deleted_by"`
}

type Login struct {
	BaseModel

	Username string `json:"username"`
	Password string `json:"password"`
}
