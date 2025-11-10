package main

import (
	"context"
	"log"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func createTemplate(ctx context.Context) prompt.ChatTemplate {
	// 创建模板，使用 FString 格式
	return prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("你是一个{role}。你需要用{style}的语气回答问题。你的目标是帮助程序员保持积极乐观的心态，提供技术建议的同时也要关注他们的心理健康。"),

		// 插入需要的对话历史（新对话的话这里不填）
		schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		schema.UserMessage("问题: {question}"),
	)
}

func createMessagesFromTemplate(ctx context.Context, convID, question string) (messages []*schema.Message, err error) {
	template := createTemplate(ctx)
	/* add start */
	chatHistory, err := eh.GetHistory(convID, 100)
	if err != nil {
		return
	}
	// 插入一条用户数据
	err = eh.SaveMessage(&schema.Message{
		Role:    schema.User,
		Content: question,
	}, convID)
	if err != nil {
		return
	}
	/* add end */
	// 使用模板生成消息
	messages, err = template.Format(context.Background(), map[string]any{
		"role":         "程序员鼓励师",
		"style":        "积极、温暖且专业",
		"question":     question, // "我的代码一直报错，感觉好沮丧，该怎么办？",
		"chat_history": chatHistory,
	})
	if err != nil {
		return
	}
	return
}

func generate(ctx context.Context, llm model.ChatModel, in []*schema.Message) *schema.Message {
	result, err := llm.Generate(ctx, in)
	if err != nil {
		log.Fatalf("llm generate failed: %v", err)
	}
	return result
}
