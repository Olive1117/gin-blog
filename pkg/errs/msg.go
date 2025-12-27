package errs

import "errors"

const (
	SUCCESS                        = "ok"
	ERROR                          = "fail"
	INVALID_PARAMS                 = "请求参数错误"
	ERROR_EXIST_TAG                = "已存在该标签名称"
	ERROR_NOT_EXIST_TAG            = "该标签不存在"
	ERROR_NOT_EXIST_ARTICLE        = "该文章不存在"
	ERROR_EXIST_CATEGORY           = "已存在该分类名称"
	ERROR_NOT_EXIST_CATEGORY       = "该分类不存在"
	ERROR_NOT_EXIST_USER           = "该用户不存在"
	ERROR_AUTH_CHECK_TOKEN_FAIL    = "Token鉴权失败"
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = "Token已超时"
	ERROR_AUTH_TOKEN               = "Token生成失败"
	ERROR_AUTH                     = "Token错误"
)

var (
	ErrNotExistUser = errors.New("该用户不存在")
	ErrAuthToken    = errors.New("Token生成失败")
)
