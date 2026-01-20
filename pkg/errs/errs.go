package errs

import "net/http"

// AppError 定义了业务错误的统一接口
type AppError struct {
	Code     int    // 业务错误码 (如 10001)
	HttpCode int    // HTTP 状态码 (如 400, 401)
	Message  string // 展示给用户的消息
	RawError error  // 原始错误（用于日志记录）
}

func (e *AppError) Error() string {
	if e.RawError != nil {
		return e.RawError.Error()
	}
	return e.Message
}

func New(httpCode, code int, msg string) *AppError {
	return &AppError{
		HttpCode: httpCode,
		Code:     code,
		Message:  msg,
	}
}

var (
	Error                    = New(http.StatusInternalServerError, 500, "fail")
	ErrInvalidParam          = New(http.StatusBadRequest, 400, "参数错误")
	ErrAuth                  = New(http.StatusUnauthorized, 10001, "登录失败")
	ErrAuthToken             = New(http.StatusInternalServerError, 10002, "令牌生成失败")
	ErrAuthCheckTokenFail    = New(http.StatusUnauthorized, 10003, "令牌错误")
	ErrAuthCheckTokenTimeout = New(http.StatusUnauthorized, 10004, "令牌已超时")
	ErrExistUser             = New(http.StatusBadRequest, 10005, "用户已存在")
	ErrNotExistUser          = New(http.StatusBadRequest, 10006, "用户不存在")
	ErrExistUsername         = New(http.StatusBadRequest, 10007, "名字已存在")
	ErrExistEmail            = New(http.StatusBadRequest, 10008, "邮箱已存在")
	ErrNotExistArticle       = New(http.StatusNotFound, 20001, "文章不存在")
	ErrNotExistCategory      = New(http.StatusNotFound, 30001, "分类不存在")
	ErrExistCategory         = New(http.StatusBadRequest, 30002, "分类已存在")
	ErrNotExistTag           = New(http.StatusNotFound, 40001, "标签不存在")
	ErrExistTag              = New(http.StatusBadRequest, 40002, "标签已存在")
	ErrNotFound              = New(http.StatusNotFound, 404, "无数据")
)
