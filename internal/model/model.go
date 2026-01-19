package model

import (
	"time"

	"github.com/Olive1117/gin-blog/pkg/idgen"
	"github.com/Olive1117/gin-blog/pkg/utils"
	"gorm.io/gorm"
)

type BaseModel struct {
	gorm.Model
	ID        int64          `gorm:"primaryKey;autoIncrement:false" json:"id,string"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedBy int64          `gorm:"column:created_by;default:NULL"`
	UpdatedBy int64          `gorm:"column:updated_by;default:NULL"`
	DeletedBy int64          `gorm:"column:deleted_by;default:NULL"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == 0 {
		b.ID = idgen.NextID()
	}
	return
}

type User struct {
	BaseModel
	// Account
	Username string `json:"username" gorm:"size:50;uniqueIndex;not null"`
	Password string `json:"-"`
}
type newUser struct {
	BaseModel

	// 账号核心
	Username string `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"` // 对应 screen_name
	Email    string `json:"email" gorm:"type:varchar(100);;not null"`
	Password string `json:"-" gorm:"type:varchar(255);not null"`

	// 基本资料
	Nickname  string `json:"nickname" gorm:"type:varchar(50)"`  // 对应 name
	Avatar    string `json:"avatar" gorm:"type:varchar(255)"`   // 对应 profile_image_url
	Banner    string `json:"banner" gorm:"type:varchar(255)"`   // 对应 profile_banner_url
	Bio       string `json:"bio" gorm:"type:text"`              // 对应 description，改为 text 类型更保险
	Location  string `json:"location" gorm:"type:varchar(100)"` // 所在地
	Website   string `json:"website" gorm:"type:varchar(255)"`  // 个人网站
	Birthdate string `json:"birthdate" gorm:"type:varchar(50)"` // 生日

	// 统计数据 (如果你想学推特做缓存计数)
	PostCount   int `json:"post_count" gorm:"default:0"`   // 对应 statuses_count
	FriendCount int `json:"friend_count" gorm:"default:0"` // 关注了多少人

	// 权限控制
	Role      string    `json:"role" gorm:"type:varchar(20);default:'user'"`
	Status    *int8     `json:"status" gorm:"type:tinyint;default:1"`
	CreatedAt time.Time `json:"created_at"`
	//TODO 最后上线时间
	// LastLoginAt time.Time `json:"last_login_at"`
}
type Article struct {
	BaseModel
	Title   string `json:"title" gorm:"size:100;not null"`
	Desc    string `json:"desc" gorm:"size:255"`
	Content string `json:"content" gorm:"type:text"`
	State   *int8  `json:"state" gorm:"default:1"`
	ShortID string `json:"short_id" gorm:"-"`

	CategoryID int64    `json:"category_id" gorm:"index"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID"`

	Tags []Tag `json:"tags" gorm:"many2many:article_tag;"`
}

func (a *Article) AfterFind(tx *gorm.DB) (err error) {
	a.ShortID = utils.EncodeByOBID(a.ID)
	return
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
