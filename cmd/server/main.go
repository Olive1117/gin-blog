package main

import (
	"fmt"
	"net/http"

	"github.com/Olive1117/gin-blog/config"
	"github.com/Olive1117/gin-blog/internal"
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
		logger.L.Error("MySQL连接失败", zap.Error(err))
	}
	gormTransaction := database.NewgormTransaction(DB)
	r := gin.New()

	internal.InitContainer(r, jwtHandler, DB, gormTransaction)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.GlobalConfig.Server.HTTPPort),
		Handler:        r,
		ReadTimeout:    config.GlobalConfig.Server.ReadTimeout,
		WriteTimeout:   config.GlobalConfig.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
