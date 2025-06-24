package gorm

import (
	"gorm.io/gorm"
)

// AutoMigrate 自动迁移所有GORM模型
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&KnowledgeBase{},
		&KnowledgeDocuments{},
		&KnowledgeChunks{},
	)
}
