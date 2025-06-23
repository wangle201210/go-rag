package gorm

import (
	"github.com/wangle201210/chat-history/models"
	"gorm.io/gorm"
)

// AutoMigrate 自动迁移所有GORM模型
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Conversation{},
		&models.Message{},
		&models.Attachment{},
		&models.MessageAttachment{},

		&KnowledgeBase{},
		&KnowledgeDocuments{},
		&KnowledgeChunks{},
	)
}
