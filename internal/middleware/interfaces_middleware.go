package middleware

import "github.com/gin-gonic/gin"

type MiddlewareContainer struct {
	Jwt         gin.HandlerFunc
	Logger      gin.HandlerFunc
	GinRecovery gin.HandlerFunc
}
