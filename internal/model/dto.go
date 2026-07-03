package model

// LoginRequest 登录请求
type LoginRequest struct {
	// 使用 binding 标签进行初步参数校验
	Username string `json:"username" binding:"required,min=1,max=32"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	// 可选：验证码标识
	// CaptchaID   string `json:"captcha_id" binding:"omitempty"`
	// CaptchaCode string `json:"captcha_code" binding:"omitempty"`
}

type ChangePasswordRequest struct {
	Username    string `json:"username" binding:"required,min=1,max=32"`
	OldPassword string `json:"old_password" binding:"required,min=6,max=64"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=64"`
}

type ArticleDTO struct {
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
	Status  *int8  `json:"status"`
	Slug    string `json:"slug"`

	CategoryName string `json:"category"`

	TagNames []string `json:"tags"`

	ImageCount int `json:"image_count"`
}

type ArticleQuery struct {
	Title        string   `form:"title"`
	CategoryName string   `form:"category"`
	TagNames     []string `form:"tags"`
	Status       *int8    `form:"status"`
}

type UserDTO struct {
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

type FriendLinkDTO struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
	Status      *int8  `json:"status"`
}
