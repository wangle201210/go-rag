package repositories

import (
	"github.com/wangle201210/go-rag/server/chat-history/models"
	"gorm.io/gorm"
)

// MessageAttachmentRepository 消息附件关联仓库
type MessageAttachmentRepository struct {
	db *gorm.DB
}

// NewMessageAttachmentRepository 创建消息附件关联仓库实例
func NewMessageAttachmentRepository(db *gorm.DB) *MessageAttachmentRepository {
	return &MessageAttachmentRepository{db: db}
}

// Create 创建消息附件关联
func (r *MessageAttachmentRepository) Create(msgAttach *models.MessageAttachment) error {
	return r.db.Create(msgAttach).Error
}

// Delete 删除消息附件关联
func (r *MessageAttachmentRepository) Delete(messageID, attachmentID string) error {
	return r.db.Where("message_id = ? AND attachment_id = ?", messageID, attachmentID).Delete(&models.MessageAttachment{}).Error
}

// ListByMessage 获取消息的所有附件关联
func (r *MessageAttachmentRepository) ListByMessage(messageID string) ([]*models.MessageAttachment, error) {
	var msgAttachs []*models.MessageAttachment
	err := r.db.Where("message_id = ?", messageID).Find(&msgAttachs).Error
	return msgAttachs, err
}

// ListByAttachment 获取附件的所有消息关联
func (r *MessageAttachmentRepository) ListByAttachment(attachmentID string) ([]*models.MessageAttachment, error) {
	var msgAttachs []*models.MessageAttachment
	err := r.db.Where("attachment_id = ?", attachmentID).Find(&msgAttachs).Error
	return msgAttachs, err
}
