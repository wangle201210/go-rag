package eino

import (
	"github.com/cloudwego/eino/schema"
	"github.com/wangle201210/go-rag/server/chat-history/models"
	"github.com/wangle201210/go-rag/server/chat-history/repositories"
)

type History struct {
	mr *repositories.MessageRepository
	cr *repositories.ConversationRepository
}

// NewEinoHistory 使用 DSN 创建 History 实例
// DSN 格式：
// - MySQL: user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
// - SQLite: /path/to/file.db 或 /path/to/file.db?_journal_mode=WAL
func NewEinoHistory(dsn string) *History {
	if err := repositories.InitDB(dsn); err != nil {
		panic(err)
	}
	return &History{
		mr: repositories.NewMessageRepository(repositories.GetDB()),
		cr: repositories.NewConversationRepository(repositories.GetDB()),
	}
}

// SaveMessage 存储message
func (x *History) SaveMessage(mess *schema.Message, convID string) error {
	return x.mr.Create(&models.Message{
		Role:           string(mess.Role),
		Content:        mess.Content,
		ConversationID: convID,
	})
}

// GetHistory 根据convID获取聊天历史
func (x *History) GetHistory(convID string, limit int) (list []*schema.Message, err error) {
	if limit == 0 {
		limit = 100
	}
	// 如果convID数据不存在，则创建
	_, err = x.cr.FirstOrCreat(convID)
	if err != nil {
		return
	}
	// 最多取100条
	mess, err := x.mr.ListByConversation(convID, 0, limit)
	if err != nil {
		return
	}
	list = messageList2ChatHistory(mess)
	return
}
