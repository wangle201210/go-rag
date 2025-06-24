package gorm

import (
	"time"
)

// KnowledgeChunks GORM模型定义
type KnowledgeChunks struct {
	ID             int64     `gorm:"primaryKey;column:id;autoIncrement"`
	KnowledgeDocID int64     `gorm:"primaryKey;column:knowledge_doc_id;not null;index"`
	ChunkID        string    `gorm:"column:chunk_id;type:varchar(36);not null;uniqueIndex:uk_chunk_id"`
	Content        string    `gorm:"column:content;type:text"`
	Ext            string    `gorm:"column:ext;type:varchar(1024)"`
	Status         int8      `gorm:"column:status;type:tinyint(1);not null;default:1"`
	CreateTime     time.Time `gorm:"column:created_at;type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)"`
	UpdateTime     time.Time `gorm:"column:updated_at;type:datetime(3);default:CURRENT_TIMESTAMP(3) on update CURRENT_TIMESTAMP(3)"`

	KnowledgeDocument KnowledgeDocuments `gorm:"foreignKey:KnowledgeDocID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:RESTRICT"`
}

// TableName 设置表名
func (KnowledgeChunks) TableName() string {
	return "knowledge_chunks"
}
