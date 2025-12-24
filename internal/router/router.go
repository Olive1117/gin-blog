package router

import "github.com/gin-gonic/gin"

func InitRouter() *gin.Engine {
	gin.SetMode("debug")
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"code": 200,
			"msg":  "hello gin!",
			"data": "",
		})
	})
	return r
}
