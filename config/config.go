package config

import (
	"os"
	"strconv"
	"time"

	"github.com/Olive1117/gin-blog/pkg/database"
)

const (
	defaultRunMode      = "debug"
	defaultHTTPPort     = 8080
	defaultReadTimeout  = 60 * time.Second
	defaultWriteTimeout = 60 * time.Second
	defaultLogLevel     = "debug"
	defaultLogPath      = "./logs/blog.log"

	defaultDBType        = "mysql"
	defaultDBUser        = "root"
	defaultDBHost        = "127.0.0.1"
	defaultDBPort        = "3306"
	defaultDBName        = "blog"
	defaultDBPassword    = ""
	defaultDBCharset     = "utf8mb4"
	defaultDBTablePrefix = "blog_"
	defaultDBLogLevel    = 4
)

var GlobalConfig = &AllConfig{
	Base: &BaseConfig{
		RunMode: "debug",
	},
	App: &AppConfig{},
	Server: &ServerConfig{
		HTTPPort:     defaultHTTPPort,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		LogLevel:     defaultLogLevel,
		LogPath:      defaultLogPath,
	},
	MySQL: &database.DBConfig{
		Host:         defaultDBHost,
		Port:         defaultDBPort,
		User:         defaultDBUser,
		DBName:       defaultDBName,
		TablePrefix:  defaultDBTablePrefix,
		Password:     defaultDBPassword,
		Charset:      defaultDBCharset,
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		MaxLifeTime:  time.Hour,
		LogLevel:     defaultDBLogLevel,
	},
}

type AllConfig struct {
	Base   *BaseConfig
	App    *AppConfig
	Server *ServerConfig
	MySQL  *database.DBConfig
}

type BaseConfig struct {
	RunMode string
}
type AppConfig struct {
	JwtSecret string
	JwtIssuer string
}
type ServerConfig struct {
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	LogPath      string
	LogLevel     string
}

func init() {
	// Base
	if v := os.Getenv("RUN_MODE"); v != "" {
		GlobalConfig.Base.RunMode = v
	}

	// App
	if v := os.Getenv("JWT_SECRET"); v != "" {
		GlobalConfig.App.JwtSecret = v
	}
	if v := os.Getenv("JWT_ISSUER"); v != "" {
		GlobalConfig.App.JwtIssuer = v
	}

	// Server
	if v := os.Getenv("HTTP_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			GlobalConfig.Server.HTTPPort = port
		}
	}
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		GlobalConfig.Server.LogLevel = v
	}

	// Database
	if v := os.Getenv("DB_HOST"); v != "" {
		GlobalConfig.MySQL.Host = v
	}
	if v := os.Getenv("DB_USER"); v != "" {
		GlobalConfig.MySQL.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		GlobalConfig.MySQL.Password = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		GlobalConfig.MySQL.DBName = v
	}
	if v := os.Getenv("DB_CHARSET"); v != "" {
		GlobalConfig.MySQL.Charset = v
	}
	if v := os.Getenv("DB_TABLE_PREFIX"); v != "" {
		GlobalConfig.MySQL.TablePrefix = v
	}
	if v := os.Getenv("DB_LOG_LEVEL"); v != "" {
		if level, err := strconv.Atoi(v); err == nil {
			GlobalConfig.MySQL.LogLevel = level
		}
	}
}
