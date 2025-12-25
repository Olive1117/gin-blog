package main

import (
	"fmt"
	"net/http"

	"github.com/Olive1117/gin-blog/config"
	"github.com/Olive1117/gin-blog/internal/router"
	"github.com/Olive1117/gin-blog/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println(config.GlobalConfig.App)
	fmt.Println(config.GlobalConfig.Base)
	fmt.Println(config.GlobalConfig.Server)
	fmt.Println(config.GlobalConfig.SQL)

	gin.SetMode(config.GlobalConfig.Base.RunMode)
	jwtHandler := jwt.NewJWT(config.GlobalConfig.App.JwtSecret, config.GlobalConfig.App.JwtIssuer)
	r := gin.New()
	router.InitRouter(r, jwtHandler)
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.GlobalConfig.Server.HTTPPort),
		Handler:        r,
		ReadTimeout:    config.GlobalConfig.Server.ReadTimeout,
		WriteTimeout:   config.GlobalConfig.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
