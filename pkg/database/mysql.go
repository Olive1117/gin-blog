package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DBConfig struct {
	Host         string        `ini:"host"`
	User         string        `ini:"user"`
	Password     string        `ini:"password"`
	DBName       string        `ini:"name"`
	TablePrefix  string        `ini:"table_prefix"`
	Charset      string        `ini:"charset"`
	MaxIdleConns int           `ini:"max_idle_conns"`
	MaxOpenConns int           `ini:"max_open_conns"`
	MaxLifeTime  time.Duration `ini:"max_life_time"`
	LogLevel     int           `ini:"log_level"` // 1: Silent, 2: Error, 3: Warn, 4: Info
}

func NewMySQLClient(cfg *DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
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
	return db, nil
}
func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}
}
