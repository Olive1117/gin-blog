package handler

import "github.com/gin-gonic/gin"

type HandlerContainer struct {
	Auth     AuthHandler
	Article  ArticleHandler
	Category CategoryHandler
	Tag      TagHandler
}
type ArticleHandler interface {
	Create(c *gin.Context)
	Delete(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
}
type AuthHandler interface {
	Auth(c *gin.Context)
}
type CategoryHandler interface {
	Create(c *gin.Context)
	Delete(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
}
type TagHandler interface {
	Create(c *gin.Context)
	Delete(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
}
