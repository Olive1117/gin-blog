package main

import (
	"fmt"
	"net/http"

	"github.com/Olive1117/gin-blog/config"
	"github.com/Olive1117/gin-blog/internal/handler"
	"github.com/Olive1117/gin-blog/internal/repository"
	"github.com/Olive1117/gin-blog/internal/router"
	"github.com/Olive1117/gin-blog/internal/service"
	"github.com/Olive1117/gin-blog/pkg/database"
	"github.com/Olive1117/gin-blog/pkg/jwt"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	logger.NewLogger(config.GlobalConfig.Server.LogPath, config.GlobalConfig.Server.LogLevel)
	gin.SetMode(config.GlobalConfig.Base.RunMode)
	jwtHandler := jwt.NewJWT(config.GlobalConfig.App.JwtSecret, config.GlobalConfig.App.JwtIssuer)
	DB, err := database.NewMySQLClient(config.GlobalConfig.MySQL)
	if err != nil {
		return
	}
	loginRepo := repository.NewLoginRepo(DB)
	loginService := service.NewLoginService(loginRepo, jwtHandler)
	loginHandler := handler.NewLoginHandler(loginService)
	r := gin.New()
	router.InitRouter(r, jwtHandler, loginHandler)
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.GlobalConfig.Server.HTTPPort),
		Handler:        r,
		ReadTimeout:    config.GlobalConfig.Server.ReadTimeout,
		WriteTimeout:   config.GlobalConfig.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
