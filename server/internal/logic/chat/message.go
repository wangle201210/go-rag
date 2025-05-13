package chat

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	role = "你是一个专业的AI助手，能够根据提供的参考信息准确回答用户问题。"
)

// createTemplate 创建并返回一个配置好的聊天模板
func createTemplate() prompt.ChatTemplate {
	// 创建模板，使用 FString 格式
	return prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("{role}"+
			"请严格遵守以下规则：\n"+
			"1. 回答必须基于提供的参考内容，不要依赖外部知识\n"+
			"2. 如果参考内容中有明确答案，直接使用参考内容回答\n"+
			"3. 如果参考内容不完整或模糊，可以合理推断但需说明\n"+
			"4. 如果参考内容完全不相关或不存在，如实告知用户'根据现有资料无法回答'\n"+
			"5. 保持回答专业、简洁、准确\n"+
			"6. 必要时可引用参考内容中的具体数据或原文\n\n"+
			"当前提供的参考内容：\n"+
			"{docs}\n\n"+
			""),
		schema.MessagesPlaceholder("chat_history", true),
		// 用户消息模板
		schema.UserMessage("Question: {question}"),
	)
}

// formatMessages 格式化消息并处理错误
func formatMessages(template prompt.ChatTemplate, data map[string]any) ([]*schema.Message, error) {
	messages, err := template.Format(context.Background(), data)
	if err != nil {
		return nil, fmt.Errorf("格式化模板失败: %w", err)
	}
	return messages, nil
}

// docsMessages 将检索到的上下文和问题转换为消息列表
func (x *Chat) docsMessages(ctx context.Context, convID string, docs []*schema.Document, question string) (messages []*schema.Message, err error) {
	chatHistory, err := x.eh.GetHistory(convID, 100)
	if err != nil {
		return
	}
	// 插入一条用户数据
	err = x.eh.SaveMessage(&schema.Message{
		Role:    schema.User,
		Content: question,
	}, convID)
	if err != nil {
		return
	}
	template := createTemplate()
	for i, doc := range docs {
		g.Log().Debugf(context.Background(), "docs[%d]: %s", i, doc.Content)
	}
	data := map[string]any{
		"role":         role,
		"question":     question,
		"docs":         docs,
		"chat_history": chatHistory,
	}
	messages, err = formatMessages(template, data)
	if err != nil {
		return
	}
	return
}
