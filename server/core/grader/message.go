package grader

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

var system = "您是一名评估检索到的文档是否足够回答用户问题的专家。" +
	"如果检索到的文档足够回答用户问题，请给出 'yes'，" +
	"如果检索到的文档不足以回答用户问题，请给出 'no'。" +
	"不要给出任何其他解释。"

// createTemplate 创建并返回一个配置好的聊天模板
func createTemplate() prompt.ChatTemplate {
	// 创建模板，使用 FString 格式
	return prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("{system}"),
		// 用户消息模板
		schema.UserMessage(
			"这是检索到的文档: \n\n"+
				" {document} \n\n"+
				"这是用户的问题: {question}"),
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

func retrieverMessages(docs []*schema.Document, question string) ([]*schema.Message, error) {
	document := ""
	for i, doc := range docs {
		document += fmt.Sprintf("docs[%d]: %s", i, doc.Content)
	}
	template := createTemplate()
	data := map[string]any{
		"system":   system,
		"question": question,
		"document": document,
	}
	messages, err := formatMessages(template, data)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
