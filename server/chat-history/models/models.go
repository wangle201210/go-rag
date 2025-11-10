package models

import (
	"encoding/json"
)

// Conversation 对话表
type Conversation struct {
	ID         uint64          `gorm:"primaryKey;column:id"`
	ConvID     string          `gorm:"uniqueIndex;column:conv_id;type:varchar(255)"`
	Title      string          `gorm:"column:title;type:varchar(255)"`
	CreatedAt  int64           `gorm:"column:created_at"`
	UpdatedAt  int64           `gorm:"column:updated_at"`
	Settings   json.RawMessage `gorm:"column:settings;type:json"`
	IsArchived bool            `gorm:"column:is_archived;default:0"`
	IsPinned   bool            `gorm:"column:is_pinned;default:0"`
}

// TableName 设置表名
func (Conversation) TableName() string {
	return "conversations"
}

// Message 消息表
type Message struct {
	ID             uint64          `gorm:"primaryKey;column:id"`
	MsgID          string          `gorm:"uniqueIndex;column:msg_id;type:varchar(255)"`
	ConversationID string          `gorm:"column:conversation_id;type:varchar(255)"`
	ParentID       string          `gorm:"column:parent_id;type:varchar(255);default:''"`
	Role           string          `gorm:"column:role;type:varchar(50);default:'user'"`
	Content        string          `gorm:"column:content;type:text"`
	CreatedAt      int64           `gorm:"column:created_at"`
	OrderSeq       int             `gorm:"column:order_seq;default:0"`
	TokenCount     int             `gorm:"column:token_count;default:0"`
	Status         string          `gorm:"column:status;type:varchar(20);default:'sent'"`
	Metadata       json.RawMessage `gorm:"column:metadata;type:json"`
	IsContextEdge  bool            `gorm:"column:is_context_edge;default:0"`
	IsVariant      bool            `gorm:"column:is_variant;default:0"`
}

// TableName 设置表名
func (Message) TableName() string {
	return "messages"
}

// Attachment 附件表
type Attachment struct {
	ID             uint64 `gorm:"primaryKey;column:id"`
	AttachID       string `gorm:"uniqueIndex;column:attach_id;type:varchar(255)"`
	MessageID      string `gorm:"column:message_id;type:varchar(255)"`
	AttachmentType string `gorm:"column:attachment_type;type:varchar(20)"`
	FileName       string `gorm:"column:file_name;type:varchar(255)"`
	FileSize       int64  `gorm:"column:file_size"`
	StorageType    string `gorm:"column:storage_type;type:varchar(20)"`
	StoragePath    string `gorm:"column:storage_path;type:varchar(1024)"`
	Thumbnail      []byte `gorm:"column:thumbnail;type:mediumblob"`
	Vectorized     bool   `gorm:"column:vectorized;default:0"`
	DataSummary    string `gorm:"column:data_summary;type:text"`
	MimeType       string `gorm:"column:mime_type;type:varchar(255)"`
	CreatedAt      int64  `gorm:"column:created_at"`
}

// TableName 设置表名
func (Attachment) TableName() string {
	return "attachments"
}

// MessageAttachment 消息附件关联表
type MessageAttachment struct {
	ID           uint64 `gorm:"primaryKey;column:id"`
	MessageID    string `gorm:"column:message_id;type:varchar(255)"`
	AttachmentID string `gorm:"column:attachment_id;type:varchar(255)"`
}

// TableName 设置表名
func (MessageAttachment) TableName() string {
	return "message_attachments"
}
