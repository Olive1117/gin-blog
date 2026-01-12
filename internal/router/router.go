package router

import (
	"github.com/Olive1117/gin-blog/internal/handler"
	"github.com/Olive1117/gin-blog/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine, handlerContainer *handler.HandlerContainer, middlewareContainer *middleware.MiddlewareContainer) {
	router.Use(gin.Logger())
	router.Use(middlewareContainer.GinRecovery)
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
	router.Use(middlewareContainer.Logger)
	public := router.Group("/api/v1")
	{
		public.GET("/test1", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"code": 200,
				"msg":  "hello gin!",
				"data": "",
			})
		})
		public.POST("/login", handlerContainer.Login.Login)
	}
	private := router.Group("/api/v1").Use(middlewareContainer.Jwt)
	{
		private.GET("/test2", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"code": 200,
				"msg":  "hello gin!",
				"data": "",
			})
		})
		private.GET("/article/:id", handlerContainer.Article.Get)
		private.POST("/article", handlerContainer.Article.Create)
		private.DELETE("/article/:id", handlerContainer.Article.Delete)
		private.PUT("/article/:id", handlerContainer.Article.Update)
		private.GET("/articles", handlerContainer.Article.List)
	}
}
