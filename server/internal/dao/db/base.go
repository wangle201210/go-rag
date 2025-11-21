package db

import (
	"fmt"

	gormModel "github.com/wangle201210/go-rag/server/internal/model/gorm"
	"gorm.io/gorm"
)

type Base struct {
	*gorm.DB
}

func (b *Base) AutoMigrate() error {
	return b.DB.AutoMigrate(gormModel.AllTables...)
}

// Ping 健康检查
func (b *Base) Ping() error {
	if b.DB == nil {
		return fmt.Errorf("database not connected")
	}

	sqlDB, err := b.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}
