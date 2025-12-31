package model

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
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"` // 过期时间戳
	TokenType   string `json:"token_type"` // 默认 "Bearer"
}

type ArticleDTO struct {
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
	State   *int8  `json:"state"`

	Category string `json:"category"`

	Tags []string `json:"tags"`
}
