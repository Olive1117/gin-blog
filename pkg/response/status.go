package response

import "net/http"

var HttpStatus = map[string]int{
	SUCCESS:                        http.StatusOK,
	ERROR:                          http.StatusInternalServerError,
	INVALID_PARAMS:                 http.StatusBadRequest,
	ERROR_EXIST_TAG:                http.StatusBadRequest,
	ERROR_NOT_EXIST_TAG:            http.StatusNotFound,
	ERROR_NOT_EXIST_ARTICLE:        http.StatusNotFound,
	ERROR_EXIST_CATEGORY:           http.StatusBadRequest,
	ERROR_NOT_EXIST_CATEGORY:       http.StatusNotFound,
	ERROR_AUTH_CHECK_TOKEN_FAIL:    http.StatusUnauthorized,
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: http.StatusUnauthorized,
	ERROR_AUTH_TOKEN:               http.StatusInternalServerError,
	ERROR_AUTH:                     http.StatusUnauthorized,
}

func GetHttpStatus(msg string) int {
	status, ok := HttpStatus[msg]
	if ok {
		return status
	}
	return HttpStatus[ERROR]
}
