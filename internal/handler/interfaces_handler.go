package handler

import "github.com/gin-gonic/gin"

type HandlerContainer struct {
	Auth     AuthHandler
	Article  ArticleHandler
	Category CategoryHandler
	Tag      TagHandler
	User     UserHandler
}
type BaseHandler interface {
	Get(c *gin.Context)
	List(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}
type ArticleHandler interface {
	BaseHandler
	Stats(c *gin.Context)
}
type AuthHandler interface {
	Auth(c *gin.Context)
}
type CategoryHandler interface {
	BaseHandler
}
type TagHandler interface {
	BaseHandler
}
type UserHandler interface {
	BaseHandler
	GetMe(c *gin.Context)
}
