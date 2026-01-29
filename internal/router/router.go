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
		articles := public.Group("/articles")
		{
			articles.GET("", handlerContainer.Article.List)
			articles.GET("/:id", handlerContainer.Article.Get)
		}
		categories := public.Group("/categories")
		{
			categories.GET("/:id", handlerContainer.Category.Get)
			categories.GET("", handlerContainer.Category.List)
		}
		tags := public.Group("/tags")
		{
			tags.GET("/:id", handlerContainer.Tag.Get)
			tags.GET("", handlerContainer.Tag.List)
		}
	}

	private := router.Group("/api/v1")
	private.Use(middlewareContainer.Jwt)
	{
		private.GET("/test2", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"code": 200,
				"msg":  "hello gin!",
				"data": "",
			})
		})
		articles := private.Group("/articles")
		{
			// articles.GET("", handlerContainer.Article.List)
			articles.POST("", handlerContainer.Article.Create)
			// articles.GET("/:id", handlerContainer.Article.Get)
			articles.PUT("/:id", handlerContainer.Article.Update)
			articles.DELETE("/:id", handlerContainer.Article.Delete)
		}

		categories := private.Group("/categories")
		{
			categories.POST("", handlerContainer.Category.Create)
			categories.DELETE("/:id", handlerContainer.Category.Delete)
			categories.PUT("/:id", handlerContainer.Category.Update)
			// categories.GET("/:id", handlerContainer.Category.Get)
			// categories.GET("", handlerContainer.Category.List)
		}

		tags := private.Group("/tags")
		{
			tags.POST("", handlerContainer.Tag.Create)
			tags.DELETE("/:id", handlerContainer.Tag.Delete)
			tags.PUT("/:id", handlerContainer.Tag.Update)
			// tags.GET("/:id", handlerContainer.Tag.Get)
			// tags.GET("", handlerContainer.Tag.List)
		}

		users := private.Group("/users")
		{
			users.POST("", handlerContainer.User.Create)
			users.DELETE("/:id", handlerContainer.User.Delete)
			users.PUT("/:id", handlerContainer.User.Update)
			users.GET("/:id", handlerContainer.User.Get)
			users.GET("", handlerContainer.User.List)
			users.GET("/me", handlerContainer.User.GetMe)
		}
	}
}
