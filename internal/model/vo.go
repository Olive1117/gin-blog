package model

import "time"

type ArticleVO struct {
	ID        int64     `json:"id,string"`
	Title     string    `json:"title"`
	Desc      string    `json:"desc"`
	Content   string    `json:"content"`
	State     *int8     `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ShortID   string    `json:"short_id" gorm:"-"`
	Slug      string    `json:"slug"`

	CategoryName string `json:"category"`

	TagNames []string `json:"tags"`

	WordCount  int `json:"word_count"`
	ImageCount int `json:"image_count"`
}

// AuthResponse 登录成功响应
type AuthResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"` // 过期时间戳
	TokenType   string    `json:"token_type"` // 默认 "Bearer"
}

type UserVO struct {
	// 账号核心
	ID       int64  `json:"id,string"`
	Username string `json:"username"`
	Email    string `json:"email"`

	// 基本资料
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	Banner    string    `json:"banner"`
	Bio       string    `json:"bio"`
	Location  string    `json:"location"`
	Website   string    `json:"website"`
	Birthdate time.Time `json:"birthdate"`

	// 统计数据 (如果你想学推特做缓存计数)
	PostCount   int `json:"post_count"`
	FriendCount int `json:"friend_count"`
	// 权限控制
	Role  string `json:"role"`
	State *int8  `json:"state"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ArticleStatsVO struct {
	Total           int64            `json:"total"`
	TotalByCategory map[string]int64 `json:"total_by_category"`
	TotalByTag      map[string]int64 `json:"total_by_tag"`
}
type CategoryVO struct {
	ID    int64  `json:"id,string"`
	Name  string `json:"name"`
	State *int8  `json:"state"`
}
type TagVO struct {
	ID    int64  `json:"id,string"`
	Name  string `json:"name"`
	State *int8  `json:"state"`
}
