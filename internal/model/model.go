package model

import (
	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model

	CreatedBy uint `gorm:"column:created_by;default:NULL"`
	UpdatedBy uint `gorm:"column:updated_by;default:NULL"`
	DeletedBy uint `gorm:"column:deleted_by;default:NULL"`
}

type User struct {
	BaseModel
	// Account
	Username string `json:"username" gorm:"size:50;uniqueIndex;not null"`
	Password string `json:"-"`
}

type Article struct {
	BaseModel
	Title   string `json:"title" gorm:"size:100;not null"`
	Desc    string `json:"desc" gorm:"size:255"`
	Content string `json:"content" gorm:"type:text"`
	State   *int8  `json:"state" gorm:"default:1"`

	CategoryID uint     `json:"category_id" gorm:"index"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID"`

	Tags []Tag `json:"tags" gorm:"many2many:article_tag;"`
}

type Category struct {
	BaseModel
	Name  string `json:"name" gorm:"size:50;uniqueIndex;not null"`
	State *int8  `json:"state" gorm:"default:1"`
}

type Tag struct {
	BaseModel
	Name  string `json:"name" gorm:"size:50;uniqueIndex;not null"`
	State *int8  `json:"state" gorm:"default:1"`

	Articles []Article `gorm:"many2many:article_tag;"`
}
