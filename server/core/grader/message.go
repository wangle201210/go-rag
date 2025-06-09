package grader

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

// createRetrieverTemplate 判断检索到的文档是否足够回答用户问题
func createRetrieverTemplate() prompt.ChatTemplate {
	// 创建模板，使用 FString 格式
	return prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage(
			"您是一名评估检索到的文档是否足够回答用户问题的专家。"+
				"请先仔细理解用户问题"+
				"如果检索到的文档足够回答用户问题，请给出 'yes'，"+
				"如果检索到的文档不足以回答用户问题，请给出 'no'。"+
				"不要给出任何其他解释。",
		),
		// 用户消息模板
		schema.UserMessage(
			"这是检索到的文档: \n"+
				"{document} \n\n"+
				"这是用户的问题: {question}"),
	)
}

// createDocRelatedTemplate 判断检索到的文档是否和用户问题相关
func createDocRelatedTemplate() prompt.ChatTemplate {
	// 创建模板，使用 FString 格式
	return prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage(
			"您是一名评估检索到的文档是否和用户问题相关的专家。"+
				"这里不需要是一个严格的测试,目标是过滤掉错误的检索。"+
				"如果检索到的文档和用户问题相关，请给出 'yes'，"+
				"如果检索到的文档和用户问题不相关，请给出 'no'。"+
				"不要给出任何其他解释。",
		),
		// 用户消息模板
		schema.UserMessage(
			"<|start_documents|> \n"+
				"{document} <|end_documents|>\n"+
				"<|start_query|>{question}<|end_query|>"),
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
	template := createRetrieverTemplate()
	data := map[string]any{
		"question": question,
		"document": document,
	}
	messages, err := formatMessages(template, data)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func docRelatedMessages(doc *schema.Document, question string) ([]*schema.Message, error) {
	template := createDocRelatedTemplate()
	data := map[string]any{
		"question": question,
		"document": doc,
	}
	messages, err := formatMessages(template, data)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
