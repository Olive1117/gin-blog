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
		public.POST("/login", handlerContainer.Auth.Auth)
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

		private.POST("/category", handlerContainer.Category.Create)
		private.DELETE("/category/:id", handlerContainer.Category.Delete)
		private.PUT("/category/:id", handlerContainer.Category.Update)
		private.GET("/category/:id", handlerContainer.Category.Get)
		private.GET("/categories", handlerContainer.Category.List)

		private.POST("/tag", handlerContainer.Tag.Create)
		private.DELETE("/tag/:id", handlerContainer.Tag.Delete)
		private.PUT("/tag/:id", handlerContainer.Tag.Update)
		private.GET("/tag/:id", handlerContainer.Tag.Get)
		private.GET("/tags", handlerContainer.Tag.List)
	}
}
