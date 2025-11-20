package db

import (
	"fmt"
	gormModel "github.com/wangle201210/go-rag/server/internal/model/gorm"
	"gorm.io/gorm"
)

type Base struct {
	db *gorm.DB
}

func (b *Base) AutoMigrate() error {
	return b.db.AutoMigrate(gormModel.AllTables...)
}

func (b *Base) DB() *gorm.DB {
	return b.db
}

func (b *Base) Disconnect() error {
	if b.db != nil {
		sqlDB, err := b.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func (b *Base) IsConnected() bool {
	if b.db == nil {
		return false
	}

	sqlDB, err := b.db.DB()
	if err != nil {
		return false
	}

	return sqlDB.Ping() == nil
}

func (b *Base) Transaction(fc func(tx *gorm.DB) error) error {
	return b.db.Transaction(fc)
}

// HasTable 检查表是否存在
func (b *Base) HasTable(model interface{}) bool {
	return b.db.Migrator().HasTable(model)
}

// GetStats 获取数据库统计信息
func (b *Base) GetStats() map[string]interface{} {
	if b.db == nil {
		return nil
	}

	sqlDB, err := b.db.DB()
	if err != nil {
		return nil
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
	}
}

// Ping 健康检查
func (b *Base) Ping() error {
	if b.db == nil {
		return fmt.Errorf("database not connected")
	}

	sqlDB, err := b.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}
