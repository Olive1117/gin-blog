package handler

import "github.com/gin-gonic/gin"

type authStore interface {
}

type authHeader struct {
	Server authStore
}

func NewAuthHeader(store authStore) *authHeader {
	return &authHeader{Server: store}
}
func (a *authHeader) login(ctx *gin.Context) {

}
