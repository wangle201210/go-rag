package db

import "gorm.io/gorm"

// Database 数据库接口
type Database interface {
	// Connect 连接管理
	Connect() error
	// DSN 构建连接DSN
	DSN() string
	// Disconnect 断开连接
	Disconnect() error
	// IsConnected 是否连接
	IsConnected() bool
	// DB 获取DB
	DB() *gorm.DB
	// AutoMigrate 迁移表结构
	AutoMigrate() error
	// Transaction 事务管理
	Transaction(fc func(tx *gorm.DB) error) error
	// HasTable 表是否存在
	HasTable(model interface{}) bool
	// GetStats 统计信息
	GetStats() map[string]interface{}
	// Ping 健康检查
	Ping() error
}

type Config struct {
	Host        string
	Port        string
	User        string
	Password    string
	Database    string
	Charset     string
	MaxOpenConn int
	MaxIdleConn int
	LogLevel    int // 1: silent 2: Error 3:Warn 4: Info

	// sqlite相关配置
	FilePath    string
	BusyTimeout int
	JournalNode string
	Synchronous string
	CacheSize   int
}
