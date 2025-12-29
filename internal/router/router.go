package router

import (
	"github.com/Olive1117/gin-blog/internal/handler"
	"github.com/Olive1117/gin-blog/internal/middleware"
	"github.com/Olive1117/gin-blog/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine, j *jwt.JWTHandler, login *handler.LoginHandler) {
	router.Use(gin.Logger())
	router.Use(middleware.GinRecovery(false))
	router.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"code": 200,
			"msg":  "hello gin!",
			"data": "",
		})
	})
	router.GET("/panic", func(c *gin.Context) {
		panic("测试：这是一个模拟的崩溃")
	})
	router.Use(middleware.GinLogger())
	public := router.Group("/api/v1")
	{
		public.GET("/test1", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"code": 200,
				"msg":  "hello gin!",
				"data": "",
			})
		})
		public.GET("/login", login.Login)
	}
	private := router.Group("/api/v1").Use(middleware.JwtAuth(j))
	{
		private.GET("/test2", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"code": 200,
				"msg":  "hello gin!",
				"data": "",
			})
		})
	}

}
