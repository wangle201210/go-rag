package repositories

import (
	"github.com/google/uuid"
	"github.com/wangle201210/go-rag/server/chat-history/models"
	"gorm.io/gorm"
)

// MessageRepository 消息仓库
type MessageRepository struct {
	db *gorm.DB
}

// NewMessageRepository 创建消息仓库实例
func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// Create 创建消息
func (r *MessageRepository) Create(msg *models.Message) error {
	if len(msg.MsgID) == 0 {
		msg.MsgID = uuid.NewString()
	}
	return r.db.Create(msg).Error
}

// Update 更新消息
func (r *MessageRepository) Update(msg *models.Message) error {
	return r.db.Save(msg).Error
}

// Delete 删除消息
func (r *MessageRepository) Delete(msgID string) error {
	return r.db.Where("msg_id = ?", msgID).Delete(&models.Message{}).Error
}

// GetByID 根据ID获取消息
func (r *MessageRepository) GetByID(msgID string) (*models.Message, error) {
	var msg models.Message
	err := r.db.Where("msg_id = ?", msgID).First(&msg).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// ListByConversation 获取对话的消息列表
func (r *MessageRepository) ListByConversation(conversationID string, offset, limit int) ([]*models.Message, error) {
	var msgs []*models.Message
	err := r.db.Where("conversation_id = ?", conversationID).
		Order("order_seq ASC").
		Offset(offset).
		Limit(limit).
		Find(&msgs).Error
	return msgs, err
}

// UpdateStatus 更新消息状态
func (r *MessageRepository) UpdateStatus(msgID string, status string) error {
	return r.db.Model(&models.Message{}).Where("msg_id = ?", msgID).Update("status", status).Error
}

// UpdateTokenCount 更新消息token数量
func (r *MessageRepository) UpdateTokenCount(msgID string, tokenCount int) error {
	return r.db.Model(&models.Message{}).Where("msg_id = ?", msgID).Update("token_count", tokenCount).Error
}

// SetContextEdge 设置消息为上下文边界
func (r *MessageRepository) SetContextEdge(msgID string, isContextEdge bool) error {
	return r.db.Model(&models.Message{}).Where("msg_id = ?", msgID).Update("is_context_edge", isContextEdge).Error
}

// SetVariant 设置消息为变体
func (r *MessageRepository) SetVariant(msgID string, isVariant bool) error {
	return r.db.Model(&models.Message{}).Where("msg_id = ?", msgID).Update("is_variant", isVariant).Error
}
