package model

import "time"

// AuthRequest 登录请求
type AuthRequest struct {
	// 使用 binding 标签进行初步参数校验
	Username string `json:"username" binding:"required,min=1,max=32"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	// 可选：验证码标识
	CaptchaID   string `json:"captcha_id" binding:"omitempty"`
	CaptchaCode string `json:"captcha_code" binding:"omitempty"`
}

// AuthResponse 登录成功响应
type AuthResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"` // 过期时间戳
	TokenType   string    `json:"token_type"` // 默认 "Bearer"
}

type ArticleDTO struct {
	ID        int64     `json:"id,string"`
	Title     string    `json:"title"`
	Desc      string    `json:"desc"`
	Content   string    `json:"content"`
	State     *int8     `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ShortID   string    `json:"short_id" gorm:"-"`

	CategoryName string `json:"category"`

	TagNames []string `json:"tags"`
}

func (dto *ArticleDTO) Category(category Category) {
	dto.CategoryName = category.Name
}
func (dto *ArticleDTO) Tags(tags []Tag) {
	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}
	dto.TagNames = tagNames
}

type ArticleQuery struct {
	Title        string   `form:"title"`
	CategoryName string   `form:"category"`
	TagNames     []string `form:"tags"`
	State        *int8    `form:"state"`
}

type CategoryDTO struct {
	ID    int64  `json:"id,string"`
	Name  string `json:"name"`
	State *int8  `json:"state"`
}
type TagDTO struct {
	ID    int64  `json:"id,string"`
	Name  string `json:"name"`
	State *int8  `json:"state"`
}

type UserDTO struct {
	// 账号核心
	ID       int64  `json:"id,string"`
	Username string `json:"username"`
	Email    string `json:"email"`

	// 基本资料
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Banner    string `json:"banner"`
	Bio       string `json:"bio"`
	Location  string `json:"location"`
	Website   string `json:"website"`
	Birthdate string `json:"birthdate"`

	// 统计数据 (如果你想学推特做缓存计数)
	PostCount   int `json:"post_count"`
	FriendCount int `json:"friend_count"`
	// 权限控制
	Role  string `json:"role"`
	State *int8  `json:"state"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type RegisterRequest struct {
	// 账号核心
	Username string `json:"username" binding:"required,min=4,max=32"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"required,min=8,max=64"`
	// 基本资料
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Banner    string `json:"banner"`
	Bio       string `json:"bio"`
	Location  string `json:"location"`
	Website   string `json:"website"`
	Birthdate string `json:"birthdate"`
}
