package model

import "time"

// LoginRequest 登录请求
type LoginRequest struct {
	// 使用 binding 标签进行初步参数校验
	Username string `json:"username" binding:"required,min=1,max=32"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	// 可选：验证码标识
	CaptchaID   string `json:"captcha_id" binding:"omitempty"`
	CaptchaCode string `json:"captcha_code" binding:"omitempty"`
}

// LoginResponse 登录成功响应
type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"` // 过期时间戳
	TokenType   string    `json:"token_type"` // 默认 "Bearer"
}

type ArticleDTO struct {
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
	State   *int8  `json:"state"`

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
	TagNames     []string `form:"tag"`
	State        *int8    `form:"state"`
}

type CategoryDTO struct {
	Name string `json:"name"`
}
type TagDTO struct {
	Name string `json:"name"`
}
