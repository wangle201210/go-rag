package gorm

import (
	"time"
)

// KnowledgeDocuments GORM模型定义
type KnowledgeDocuments struct {
	ID                int64     `gorm:"primaryKey;column:id;autoIncrement"`
	KnowledgeBaseName string    `gorm:"column:knowledge_base_name;type:varchar(255);not null"`
	FileName          string    `gorm:"column:file_name;type:varchar(255)"`
	Status            int8      `gorm:"column:status;type:tinyint;not null;default:0"`
	CreateTime        time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdateTime        time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP"`
}

// TableName 设置表名
func (KnowledgeDocuments) TableName() string {
	return "knowledge_documents"
}
