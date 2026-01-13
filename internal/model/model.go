package model

import (
	"time"

	"gorm.io/gorm"
)

// TODO 改掉uint
type BaseModel struct {
	gorm.Model
	ID        int64          `gorm:"primaryKey" json:"id,string"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy int64          `gorm:"column:created_by;default:NULL"`
	UpdatedBy int64          `gorm:"column:updated_by;default:NULL"`
	DeletedBy int64          `gorm:"column:deleted_by;default:NULL"`
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

	CategoryID int64    `json:"category_id" gorm:"index"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID"`

	Tags []Tag `json:"tags" gorm:"many2many:article_tag;"`
}

func (a *Article) CategoryName(categoryName string) {
	a.Category = Category{Name: categoryName}
}
func (a *Article) TagNames(tagNames []string) {
	tags := make([]Tag, len(tagNames))
	for i, name := range tagNames {
		tags[i].Name = name
	}
	a.Tags = tags
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
