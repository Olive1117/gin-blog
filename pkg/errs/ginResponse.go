package errs

import (
	"github.com/gin-gonic/gin"
)

type ResponseBody struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Result(c *gin.Context, httpCode int, msg string, data any) {
	if data == nil {
		data = gin.H{}
	}
	c.JSON(httpCode, ResponseBody{
		Code: GetAppCode(msg),
		Msg:  msg,
		Data: data,
	})
}

func Success(c *gin.Context, data any) {
	Result(c, HttpStatus[SUCCESS], SUCCESS, data)
}

func Error(c *gin.Context, msg string, data any) {
	Result(c, GetHttpStatus(msg), msg, data)
}
