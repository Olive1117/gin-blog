package config

import (
	"embed"
	"log"
	"time"

	"github.com/Olive1117/gin-blog/pkg/database"
	"github.com/go-ini/ini"
)

var GlobalConfig = &AllConfig{
	Base: &BaseConfig{
		RunMode: "debug",
	},
	App: &AppConfig{
		// JwtSecret: "!@)*#)!@U#@*!@!)",
	},
	Server: &ServerConfig{
		HTTPPort:     8000,
		ReadTimeout:  60,
		WriteTimeout: 60,
	},
	SQL: &SQLConfig{
		Type:        "mysql",
		User:        "root",
		Host:        "127.0.0.1:3306",
		Name:        "bolg",
		TablePrefix: "blog_",
		Charset:     "utf8mb4",
	},
	MySQL: &database.DBConfig{
		Host:         "127.0.0.1:3306",
		User:         "root",
		DBName:       "blog",
		TablePrefix:  "blog_",
		Charset:      "utf8mb4",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		MaxLifeTime:  time.Hour,
		LogLevel:     3,
	},
}

//go:embed app.ini
var ConfigFS embed.FS

type AllConfig struct {
	Base   *BaseConfig
	App    *AppConfig
	Server *ServerConfig
	SQL    *SQLConfig
	MySQL  *database.DBConfig
}

type BaseConfig struct {
	RunMode string `ini:"run_mode"`
}
type AppConfig struct {
	JwtSecret string `ini:"jwt_secret"`
	JwtIssuer string `ini:"jwt_issuer"`
}
type ServerConfig struct {
	HTTPPort     int           `ini:"http_port"`
	ReadTimeout  time.Duration `ini:"read_timeout"`
	WriteTimeout time.Duration `ini:"write_timeout"`
}
type SQLConfig struct {
	Type        string `ini:"type"`
	User        string `ini:"user"`
	Password    string `ini:"password"`
	Host        string `ini:"host"`
	Name        string `ini:"name"`
	TablePrefix string `ini:"table_prefix"`
	Charset     string `ini:"charset"`
}

func init() {
	//TODO 移除硬编码文件名
	data, err := ConfigFS.ReadFile("app.ini")
	if err != nil {
		log.Fatalf("没有找到app.ini配置文件")
	}
	cfg, err := ini.Load(data)
	if err != nil {
		log.Fatalf("app.ini读取错误")
	}
	// 映射配置
	cfg.Section("").MapTo(GlobalConfig.Base)
	cfg.Section("app").MapTo(GlobalConfig.App)
	cfg.Section("server").MapTo(GlobalConfig.Server)
	cfg.Section("database").MapTo(GlobalConfig.SQL)
	cfg.Section("database").MapTo(GlobalConfig.MySQL)
	// 处理时间转换
	GlobalConfig.Server.ReadTimeout *= time.Second
	GlobalConfig.Server.WriteTimeout *= time.Second
}
