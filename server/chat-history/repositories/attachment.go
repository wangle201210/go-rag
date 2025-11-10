package repositories

import (
	"github.com/wangle201210/go-rag/server/chat-history/models"
	"gorm.io/gorm"
)

// AttachmentRepository 附件仓库
type AttachmentRepository struct {
	db *gorm.DB
}

// NewAttachmentRepository 创建附件仓库实例
func NewAttachmentRepository(db *gorm.DB) *AttachmentRepository {
	return &AttachmentRepository{db: db}
}

// Create 创建附件
func (r *AttachmentRepository) Create(attach *models.Attachment) error {
	return r.db.Create(attach).Error
}

// Update 更新附件
func (r *AttachmentRepository) Update(attach *models.Attachment) error {
	return r.db.Save(attach).Error
}

// Delete 删除附件
func (r *AttachmentRepository) Delete(attachID string) error {
	return r.db.Where("attach_id = ?", attachID).Delete(&models.Attachment{}).Error
}

// GetByID 根据ID获取附件
func (r *AttachmentRepository) GetByID(attachID string) (*models.Attachment, error) {
	var attach models.Attachment
	err := r.db.Where("attach_id = ?", attachID).First(&attach).Error
	if err != nil {
		return nil, err
	}
	return &attach, nil
}

// ListByMessage 获取消息的附件列表
func (r *AttachmentRepository) ListByMessage(messageID string) ([]*models.Attachment, error) {
	var attachs []*models.Attachment
	err := r.db.Where("message_id = ?", messageID).Find(&attachs).Error
	return attachs, err
}

// UpdateVectorized 更新附件向量化状态
func (r *AttachmentRepository) UpdateVectorized(attachID string, vectorized bool) error {
	return r.db.Model(&models.Attachment{}).Where("attach_id = ?", attachID).Update("vectorized", vectorized).Error
}

// UpdateDataSummary 更新附件数据摘要
func (r *AttachmentRepository) UpdateDataSummary(attachID string, dataSummary string) error {
	return r.db.Model(&models.Attachment{}).Where("attach_id = ?", attachID).Update("data_summary", dataSummary).Error
}
