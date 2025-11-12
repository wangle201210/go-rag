package eino

import (
	"github.com/cloudwego/eino/schema"
	"github.com/wangle201210/go-rag/server/chat-history/models"
)

func messageList2ChatHistory(mess []*models.Message) (history []*schema.Message) {
	for _, m := range mess {
		history = append(history, message2MessagesTemplate(m))
	}
	return
}

func message2MessagesTemplate(mess *models.Message) *schema.Message {
	return &schema.Message{
		Role:    schema.RoleType(mess.Role),
		Content: mess.Content,
	}
}
