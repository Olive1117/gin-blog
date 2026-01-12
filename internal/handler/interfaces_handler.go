package handler

import "github.com/gin-gonic/gin"

type HandlerContainer struct {
	Login   LoginHandler
	Article ArticleHandler
}
type ArticleHandler interface {
	Create(c *gin.Context)
	Delete(c *gin.Context)
	Get(c *gin.Context)
	List(c *gin.Context)
	Update(c *gin.Context)
}
type LoginHandler interface {
	Login(c *gin.Context)
}
