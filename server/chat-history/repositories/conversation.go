package repositories

import (
	"github.com/wangle201210/go-rag/server/chat-history/models"
	"gorm.io/gorm"
)

// ConversationRepository 对话仓库
type ConversationRepository struct {
	db *gorm.DB
}

// NewConversationRepository 创建对话仓库实例
func NewConversationRepository(db *gorm.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

// Create 创建对话
func (r *ConversationRepository) Create(conv *models.Conversation) error {
	return r.db.Create(conv).Error
}

// Update 更新对话
func (r *ConversationRepository) Update(conv *models.Conversation) error {
	return r.db.Save(conv).Error
}

// Delete 删除对话
func (r *ConversationRepository) Delete(convID string) error {
	return r.db.Where("conv_id = ?", convID).Delete(&models.Conversation{}).Error
}

// GetByID 根据ID获取对话
func (r *ConversationRepository) GetByID(convID string) (*models.Conversation, error) {
	var conv models.Conversation
	err := r.db.Where("conv_id = ?", convID).First(&conv).Error
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

// FirstOrCreat 根据ID判断，如果存在就返回，不存在就创建
func (r *ConversationRepository) FirstOrCreat(convID string) (*models.Conversation, error) {
	var conv models.Conversation
	err := r.db.Where(models.Conversation{ConvID: convID}).FirstOrCreate(&conv).Error
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

// List 获取对话列表
func (r *ConversationRepository) List(offset, limit int) ([]*models.Conversation, error) {
	var convs []*models.Conversation
	err := r.db.Offset(offset).Limit(limit).Order("updated_at DESC").Find(&convs).Error
	return convs, err
}

// Archive 归档对话
func (r *ConversationRepository) Archive(convID string) error {
	return r.db.Model(&models.Conversation{}).Where("conv_id = ?", convID).Update("is_archived", true).Error
}

// Unarchive 取消归档对话
func (r *ConversationRepository) Unarchive(convID string) error {
	return r.db.Model(&models.Conversation{}).Where("conv_id = ?", convID).Update("is_archived", false).Error
}

// Pin 置顶对话
func (r *ConversationRepository) Pin(convID string) error {
	return r.db.Model(&models.Conversation{}).Where("conv_id = ?", convID).Update("is_pinned", true).Error
}

// Unpin 取消置顶对话
func (r *ConversationRepository) Unpin(convID string) error {
	return r.db.Model(&models.Conversation{}).Where("conv_id = ?", convID).Update("is_pinned", false).Error
}
