package main

import (
	"fmt"
	"net/http"

	"github.com/Olive1117/gin-blog/config"
	"github.com/Olive1117/gin-blog/internal/router"
)

func main() {
	fmt.Println(config.GlobalConfig.App)
	fmt.Println(config.GlobalConfig.Base)
	fmt.Println(config.GlobalConfig.Server)
	fmt.Println(config.GlobalConfig.SQL)
	router := router.InitRouter()
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.GlobalConfig.Server.HTTPPort),
		Handler:        router,
		ReadTimeout:    config.GlobalConfig.Server.ReadTimeout,
		WriteTimeout:   config.GlobalConfig.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	server.ListenAndServe()
}
