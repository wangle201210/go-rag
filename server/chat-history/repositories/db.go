package repositories

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/wangle201210/go-rag/server/chat-history/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// InitDB 初始化数据库连接
// dsn 格式：
// - MySQL: user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
// - SQLite: /path/to/file.db 或 /path/to/file.db?_journal_mode=WAL
func InitDB(dsn string) error {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 根据 DSN 格式判断数据库类型
	var dialector gorm.Dialector
	var dbType string
	var maxOpenConns, maxIdleConns int

	if strings.Contains(dsn, "@tcp(") {
		// MySQL DSN 格式
		dbType = "mysql"
		dialector = mysql.Open(dsn)
		maxOpenConns = 100
		maxIdleConns = 10
	} else {
		// SQLite DSN 格式（文件路径）
		dbType = "sqlite"
		dialector = sqlite.Open(dsn)
		maxOpenConns = 1
		maxIdleConns = 1
	}

	var err error
	db, err = gorm.Open(dialector, config)
	if err != nil {
		return fmt.Errorf("failed to connect database (%s): %v", dbType, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移数据库表结构
	if err = autoMigrateTables(); err != nil {
		return fmt.Errorf("failed to migrate database tables: %v", err)
	}

	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("database connection not initialized")
	}
	return db
}

// autoMigrateTables 自动迁移数据库表结构
func autoMigrateTables() error {
	// 自动迁移会创建表、缺失的外键、约束、列和索引
	return db.AutoMigrate(
		&models.Conversation{},
		&models.Message{},
		&models.Attachment{},
		&models.MessageAttachment{},
	)
}
