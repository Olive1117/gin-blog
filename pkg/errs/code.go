package errs

var AppCode = map[string]int{
	SUCCESS:                        200,
	ERROR:                          500,
	INVALID_PARAMS:                 400,
	ERROR_EXIST_TAG:                400,
	ERROR_NOT_EXIST_TAG:            404,
	ERROR_NOT_EXIST_ARTICLE:        404,
	ERROR_EXIST_CATEGORY:           400,
	ERROR_NOT_EXIST_CATEGORY:       404,
	ERROR_AUTH_CHECK_TOKEN_FAIL:    401,
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: 401,
	ERROR_AUTH_TOKEN:               500,
	ERROR_AUTH:                     401,
}

func GetAppCode(msg string) int {
	code, ok := AppCode[msg]
	if ok {
		return code
	}
	return AppCode[ERROR]
}
