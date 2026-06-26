package model

import (
	"time"

	"github.com/Olive1117/gin-blog/pkg/idgen"
	"github.com/Olive1117/gin-blog/pkg/utils"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int64          `gorm:"primaryKey;autoIncrement:false;type:bigint(20) unsigned" json:"id,string"`
	CreatedAt time.Time      `gorm:"type:datetime(3)" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime(3)" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:datetime(3);index:,composite:deletedat" json:"-"`
	CreatedBy int64          `gorm:"type:bigint(20) unsigned;default:0;comment:创建者ID"`
	UpdatedBy int64          `gorm:"type:bigint(20) unsigned;default:0;comment:修改者ID"`
	DeletedBy int64          `gorm:"type:bigint(20) unsigned;default:0;comment:删除者ID"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == 0 {
		b.ID = idgen.NextID()
	}
	return
}

type User struct {
	BaseModel

	// 账号核心
	Username string `json:"username" gorm:"type:varchar(50);not null;default:'';comment:账号;index:,composite:deletedat,unique,priority:1"` // 对应 screen_name
	Email    string `json:"email" gorm:"type:varchar(100);not null;default:'';comment:邮箱"`
	Password string `json:"-" gorm:"type:varchar(255);not null;default:'';comment:密码"`

	// 基本资料
	Nickname  string     `json:"nickname" gorm:"type:varchar(50);default:'';comment:昵称"`  // 对应 name
	Avatar    string     `json:"avatar" gorm:"type:varchar(255);default:'';comment:头像"`   // 对应 profile_image_url
	Banner    string     `json:"banner" gorm:"type:varchar(255);default:'';comment:背景"`   // 对应 profile_banner_url
	Bio       string     `json:"bio" gorm:"type:text;comment:个人简介"`                       // 对应 description，改为 text 类型更保险
	Location  string     `json:"location" gorm:"type:varchar(100);default:'';comment:地址"` // 所在地
	Website   string     `json:"website" gorm:"type:varchar(255);default:'';comment:网站"`  // 个人网站
	Birthdate *time.Time `json:"birthdate" gorm:"type:date;comment:生日"`                   // 生日

	// 统计数据 (如果你想学推特做缓存计数)
	PostCount   int `json:"post_count" gorm:"type:int(10) unsigned;not null;default:0;comment:文章数量"`   // 对应 statuses_count
	FriendCount int `json:"friend_count" gorm:"type:int(10) unsigned;not null;default:0;comment:好友数量"` // 关注了多少人

	// 权限控制
	Role  string `json:"role" gorm:"type:varchar(20);not null;default:user;comment:权限"`
	State *int8  `json:"state" gorm:"type:tinyint(3) unsigned;not null;default:1;comment:状态: 1-正常 2-冻结"`
	//TODO 最后上线时间
	// LastLoginAt time.Time `json:"last_login_at"`
}

func (u *User) BirthdateString(BirthdateString string) {
	if b, err := time.ParseInLocation("2006-01-02", BirthdateString, time.Local); err == nil {
		u.Birthdate = &b
	}
}

type Article struct {
	BaseModel
	Title   string `json:"title" gorm:"type:varchar(100);default:'';comment:文章标题"`
	Desc    string `json:"desc" gorm:"type:varchar(255);default:'';comment:简述"`
	Content string `json:"content" gorm:"type:text"`
	State   *int8  `json:"state" gorm:"type:tinyint(3) unsigned;default:1;comment:状态 0为禁用1为启用"`
	ShortID string `json:"short_id" gorm:"-"`

	CategoryID int64    `json:"category_id" gorm:"type:bigint(20) unsigned;default:0;comment:分类ID"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID"`

	Tags []Tag `json:"tags" gorm:"many2many:article_tag;"`

	WordCount  int `json:"word_count" gorm:"type:int(10) unsigned;default:0;comment:文章字数"`
	ImageCount int `json:"image_count" gorm:"type:int(10) unsigned;default:0;comment:文章字数"`
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
	Name  string `json:"name" gorm:"type:varchar(100);default:'';comment:分类名称;index:,composite:deletedat,unique,priority:1"`
	State *int8  `json:"state" gorm:"type:tinyint(3) unsigned;default:1;comment:状态 0为禁用、1为启用"`
}

type Tag struct {
	BaseModel
	Name  string `json:"name" gorm:"type:varchar(100);default:'';comment:标签名称;index:,composite:deletedat,unique,priority:1"`
	State *int8  `json:"state" gorm:"type:tinyint(3) unsigned;default:1;comment:状态 0为禁用、1为启用"`

	Articles []Article `gorm:"many2many:article_tag;"`
}

type ArticleTag struct {
	ArticleID uint64  `gorm:"type:bigint(20) unsigned;primaryKey;comment:文章ID" json:"article_id"`
	TagID     uint64  `gorm:"type:bigint(20) unsigned;primaryKey;comment:标签ID" json:"tag_id"`
	Article   Article `gorm:"constraint:OnDelete:CASCADE;foreignKey:ArticleID" json:"-"`
	Tag       Tag     `gorm:"constraint:OnDelete:CASCADE;foreignKey:TagID" json:"-"`
}
