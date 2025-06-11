package core

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

var system = "你非常擅长于使用rag进行数据检索，" +
	"你的目标是在充分理解用户的问题后进行向量化检索\n" +
	"现在时间{time_now}\n" +
	"你要优化并提取搜索的查询内容。" +
	"请遵循以下规则重写查询内容：\n" +
	"- 根据用户的问题和上下文，重写应该进行搜索的关键词\n" +
	"- 如果需要使用时间，则根据当前时间给出需要查询的具体时间日期信息\n" +
	// "- 生成的查询关键词要选择合适的语言，考虑用户的问题类型使用最适合的语言进行搜索，例如某些问题应该保持用户的问题语言，而有一些则更适合翻译成英语或其他语言\n" +
	"- 保持查询简洁，查询内容通常不超过3个关键词, 最多不要超过5个关键词\n" +
	"- 参考Elasticsearch搜索查询习惯重写关键字。" +
	"- 直接返回优化后的搜索词，不要有任何额外说明。\n" +
	"- 尽量不要使用下面这些已使用过的关键词，因为之前使用这些关键词搜索到的结果不符合预期，已使用过的关键词：{used}\n" +
	"- 尽量不使用知识库名字《{knowledgeBase}》中包含的关键词\n"

// createTemplate 创建并返回一个配置好的聊天模板
func createTemplate() prompt.ChatTemplate {
	return prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage(system),
		// 用户消息模板
		schema.UserMessage(
			"如下是用户的问题: {question}"),
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

func getOptimizedQueryMessages(used, question, knowledgeBase string) ([]*schema.Message, error) {
	template := createTemplate()
	data := map[string]any{
		"time_now":      time.Now().Format(time.RFC3339),
		"question":      question,
		"used":          used,
		"knowledgeBase": knowledgeBase,
	}
	messages, err := formatMessages(template, data)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
