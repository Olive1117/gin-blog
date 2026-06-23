package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Olive1117/gin-blog/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DBConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	TablePrefix  string
	Charset      string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifeTime  time.Duration
	LogLevel     int // 1: Silent, 2: Error, 3: Warn, 4: Info
}

func NewMySQLClient(cfg *DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(cfg.LogLevel)),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.TablePrefix,
			SingularTable: true,
		},
	}
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	InitSchemaAndSeed(db)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifeTime)
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}
	db.Use(&AuditPlugin{})
	return db, nil
}
func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}
}

type gormTransaction struct {
	db *gorm.DB
}

func NewgormTransaction(db *gorm.DB) *gormTransaction {
	return &gormTransaction{
		db: db,
	}
}

func (g *gormTransaction) Transaction(c context.Context, fn func(c context.Context) error) error {
	return g.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		newc := context.WithValue(c, "tx", tx)
		return fn(newc)
	})
}

func InitSchemaAndSeed(db *gorm.DB) {
	if !db.Migrator().HasTable(&model.Tag{}) {
		db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章标签管理'").AutoMigrate(&model.Tag{})
	}

	if !db.Migrator().HasTable(&model.Category{}) {
		db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章分类管理'").AutoMigrate(&model.Category{})
	}

	if !db.Migrator().HasTable(&model.Article{}) {
		db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章管理'").AutoMigrate(&model.Article{})
	}
	if !db.Migrator().HasTable(&model.ArticleTag{}) {
		db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章标签关联'").AutoMigrate(&model.ArticleTag{})
	}
	if !db.Migrator().HasTable(&model.User{}) {
		db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户管理'").AutoMigrate(&model.User{})
	}
	// 插入默认 admin 用户
	var count int64
	db.Model(&model.User{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		db.Create(&model.User{
			BaseModel: model.BaseModel{ID: 1},
			Username:  "admin",
			Password:  "123456", // 实际项目要加密
		})
	}
}
