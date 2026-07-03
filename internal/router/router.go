package router

import (
	"github.com/Olive1117/gin-blog/internal/handler"
	"github.com/Olive1117/gin-blog/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine, handlerContainer *handler.HandlerContainer, middlewareContainer *middleware.MiddlewareContainer) {
	router.Use(middlewareContainer.GinRecovery)
	router.Use(gin.Logger())
	router.Use(middlewareContainer.Logger)
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

	api := router.Group("/api/v1")

	api.GET("/test1", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"code": 200,
			"msg":  "hello gin!",
			"data": "",
		})
	})
	// 文章相关路由
	api.GET("/articles", handlerContainer.Article.List)
	api.GET("/articles/:id", handlerContainer.Article.Get)
	api.GET("/articles/stats", handlerContainer.Article.Stats)
	// 分类相关路由
	api.GET("/categories/:id", handlerContainer.Category.Get)
	api.GET("/categories", handlerContainer.Category.List)
	// 标签相关路由
	api.GET("/tags/:id", handlerContainer.Tag.Get)
	api.GET("/tags", handlerContainer.Tag.List)
	// 用户相关路由
	api.POST("/login", handlerContainer.User.Login)
	api.POST("/users", handlerContainer.User.Create)

	// 需要认证的路由组
	auth := api.Group("")
	auth.Use(middlewareContainer.Jwt)
	auth.GET("/test2", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"code": 200,
			"msg":  "hello gin!",
			"data": "",
		})
	})
	// 文章相关路由
	auth.POST("/articles", handlerContainer.Article.Create)
	auth.PUT("/articles/:id", handlerContainer.Article.Update)
	auth.DELETE("/articles/:id", handlerContainer.Article.Delete)
	// 分类相关路由
	auth.POST("/categories", handlerContainer.Category.Create)
	auth.DELETE("/categories/:id", handlerContainer.Category.Delete)
	auth.PUT("/categories/:id", handlerContainer.Category.Update)
	// 标签相关路由
	auth.POST("/tags", handlerContainer.Tag.Create)
	auth.DELETE("/tags/:id", handlerContainer.Tag.Delete)
	auth.PUT("/tags/:id", handlerContainer.Tag.Update)
	// 用户相关路由
	auth.DELETE("/users/:id", handlerContainer.User.Delete)
	auth.PUT("/users/:id", handlerContainer.User.Update)
	auth.GET("/users/:id", handlerContainer.User.Get)
	auth.GET("/users", handlerContainer.User.List)
	auth.GET("/users/me", handlerContainer.User.GetMe)
}
