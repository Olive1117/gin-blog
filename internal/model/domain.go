package model

import (
	"time"

	"github.com/Olive1117/gin-blog/pkg/utils"
	"gorm.io/gorm"
)

type User struct {
	BaseModel

	// 账号核心
	Username string `json:"username" gorm:"size:50;not null;default:'';comment:账号;index:,composite:deletedat,unique,priority:1"` // 对应 screen_name
	Email    string `json:"email" gorm:"size:100;not null;default:'';comment:邮箱"`
	Password string `json:"-" gorm:"size:255;not null;default:'';comment:密码"`

	// 基本资料
	Nickname  string     `json:"nickname" gorm:"size:50;default:'';comment:昵称"`  // 对应 name
	Avatar    string     `json:"avatar" gorm:"size:255;default:'';comment:头像"`   // 对应 profile_image_url
	Banner    string     `json:"banner" gorm:"size:255;default:'';comment:背景"`   // 对应 profile_banner_url
	Bio       string     `json:"bio" gorm:"type:text;comment:个人简介"`              // 对应 description，改为 text 类型更保险
	Location  string     `json:"location" gorm:"size:100;default:'';comment:地址"` // 所在地
	Website   string     `json:"website" gorm:"size:255;default:'';comment:网站"`  // 个人网站
	Birthdate *time.Time `json:"birthdate" gorm:"comment:生日"`                    // 生日

	// 统计数据 (如果你想学推特做缓存计数)
	PostCount   int `json:"post_count" gorm:"default:0;comment:文章数量"`   // 对应 statuses_count
	FriendCount int `json:"friend_count" gorm:"default:0;comment:好友数量"` // 关注了多少人

	// 权限控制
	Role   string `json:"role" gorm:"size:20;not null;default:'user';comment:权限"`
	Status *int8  `json:"status" gorm:"default:1;comment:状态 0-禁用 1-启用"`
	//TODO 最后上线时间
	// LastLoginAt time.Time `json:"last_login_at"`
}

type Article struct {
	BaseModel
	Title   string `json:"title" gorm:"size:100;default:'';comment:文章标题"`
	Desc    string `json:"desc" gorm:"size:255;default:'';comment:简述"`
	Content string `json:"content" gorm:"type:text"`
	ShortID string `json:"short_id" gorm:"-"`
	Slug    string `json:"slug" gorm:"size:255;not null;default:'';comment:URL尾链"`
	Status  *int8  `json:"status" gorm:"default:1;comment:状态 0-禁用 1-启用"`

	CategoryID int64    `json:"category_id" gorm:"default:0;comment:分类ID"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID"`

	Tags []Tag `json:"tags" gorm:"many2many:article_tag;"`

	WordCount  int `json:"word_count" gorm:"default:0;comment:文章字数"`
	ImageCount int `json:"image_count" gorm:"default:0;comment:文章字数"`
}

func (a *Article) AfterFind(tx *gorm.DB) (err error) {
	a.ShortID = utils.EncodeByOBID(a.ID)
	return
}

type Category struct {
	BaseModel
	Name   string `json:"name" gorm:"size:100;default:'';comment:分类名称;index:,composite:deletedat,unique,priority:1"`
	Status *int8  `json:"status" gorm:"default:1;comment:状态 0-禁用 1-启用"`
}

type Tag struct {
	BaseModel
	Name   string `json:"name" gorm:"size:100;default:'';comment:标签名称;index:,composite:deletedat,unique,priority:1"`
	Status *int8  `json:"status" gorm:"default:1;comment:状态 0-禁用 1-启用"`

	Articles []Article `gorm:"many2many:article_tag;"`
}

type ArticleTag struct {
	ArticleID int64   `gorm:"primaryKey;not null;comment:文章ID" json:"article_id"`
	TagID     int64   `gorm:"primaryKey;not null;comment:标签ID" json:"tag_id"`
	Article   Article `gorm:"constraint:OnDelete:CASCADE;foreignKey:ArticleID" json:"-"`
	Tag       Tag     `gorm:"constraint:OnDelete:CASCADE;foreignKey:TagID" json:"-"`
}

type FriendLink struct {
	BaseModel
	Name        string `json:"name" gorm:"size:100;default:'';comment:友链名称"`
	URL         string `json:"url" gorm:"size:255;default:'';comment:友链URL"`
	Logo        string `json:"logo" gorm:"size:255;default:'';comment:友链Logo"`
	Description string `json:"description" gorm:"size:255;comment:友链描述"`
	Status      *int8  `json:"status" gorm:"default:1;comment:状态 0-禁用 1-启用"`
}

type Bookmark struct {
	BaseModel

	Name   string `json:"name" gorm:"size:100;not null;comment:链接名称"`
	URL    string `json:"url" gorm:"size:255;not null;comment:链接地址"`
	Folder string `json:"folder" gorm:"size:100;default:'';comment:所属文件夹"`
	UserID int64  `json:"user_id" gorm:"not null;comment:用户ID;index"`
	Remark string `json:"remark" gorm:"size:255;default:'';comment:备注"`
}

type Diary struct {
	BaseModel

	Content   string     `json:"content" gorm:"not null;comment:日记内容"`
	Mood      string     `json:"mood" gorm:"size:20;default:'';comment:心情"`
	Weather   string     `json:"weather" gorm:"size:20;default:'';comment:天气"`
	DiaryDate *time.Time `json:"diary_date" gorm:"type:date;not null;comment:日记日期;index"`
	IsPublic  bool       `json:"is_public" gorm:"default:false;comment:是否公开"`
}
