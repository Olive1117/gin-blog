package database

import (
	"context"
	"fmt"
	"time"

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
