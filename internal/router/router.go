package router

import (
	"github.com/Olive1117/gin-blog/internal/middleware"
	"github.com/Olive1117/gin-blog/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine, j *jwt.JWTHandler) {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"code": 200,
			"msg":  "hello gin!",
			"data": "",
		})
	})

	public := router.Group("/api/v1")
	{
		public.GET("/test1", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"code": 200,
				"msg":  "hello gin!",
				"data": "",
			})
		})
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
