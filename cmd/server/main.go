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
	"go.uber.org/zap"
)

func main() {
	logger.NewLogger(config.GlobalConfig.Server.LogPath, config.GlobalConfig.Server.LogLevel)
	gin.SetMode(config.GlobalConfig.Base.RunMode)
	jwtHandler := jwt.NewJWT(config.GlobalConfig.App.JwtSecret, config.GlobalConfig.App.JwtIssuer)
	DB, err := database.NewMySQLClient(config.GlobalConfig.MySQL)
	if err != nil {
		logger.L.Error("MySQL初始化失败", zap.Error(err))
		return
	}
	gormTransaction := database.NewgormTransaction(DB)
	loginHandler := handler.NewLoginHandler(service.NewLoginService(repository.NewLoginRepo(DB), jwtHandler))
	articleHandler := handler.NewArticleHandler(service.NewArticleService(repository.NewArticleRepo(DB), gormTransaction))
	r := gin.New()
	router.InitRouter(r, jwtHandler, loginHandler, articleHandler)
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.GlobalConfig.Server.HTTPPort),
		Handler:        r,
		ReadTimeout:    config.GlobalConfig.Server.ReadTimeout,
		WriteTimeout:   config.GlobalConfig.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
