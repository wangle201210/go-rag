package gorm

import (
	"gorm.io/gorm"
)

var AllTables = []any{
	&KnowledgeBase{},
	&KnowledgeDocuments{},
	&KnowledgeChunks{},
}

// AutoMigrate 自动迁移所有GORM模型
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		AllTables...,
	)
}
