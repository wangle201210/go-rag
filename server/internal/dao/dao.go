package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	gormModel "github.com/wangle201210/go-rag/server/internal/model/gorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func init() {
	err := InitDB()
	if err != nil {
		g.Log().Fatal(context.Background(), "database connection not initialized")
	}
}

// InitDB 初始化数据库连接
func InitDB() error {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	dsn := GetDsn()

	var err error
	db, err = gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移数据库表结构
	if err = gormModel.AutoMigrate(db); err != nil {
		return fmt.Errorf("failed to migrate database tables: %v", err)
	}

	return nil
}

func GetDsn() string {
	cfg := g.DB().GetConfig()
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name)
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	if db == nil {
		g.Log().Fatal(context.Background(), "database connection not initialized")
	}
	return db
}
