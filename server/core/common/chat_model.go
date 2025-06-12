package common

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/model"
	"github.com/gogf/gf/v2/frame/g"
)

var chatModel model.BaseChatModel
var noThinkChatModel model.BaseChatModel

func GetChatModel(ctx context.Context, cfg *openai.ChatModelConfig) (model.BaseChatModel, error) {
	if chatModel != nil {
		return chatModel, nil
	}
	cm, err := openai.NewChatModel(ctx, cfg)
	if err != nil {
		return nil, err
	}
	chatModel = cm
	return cm, nil
}

func GetNotThinkChatModel(ctx context.Context, cfg *qwen.ChatModelConfig) (model.BaseChatModel, error) {
	if noThinkChatModel != nil {
		return noThinkChatModel, nil
	}
	// &config.Config{
	// 	Client:         client,
	// 	IndexName:      "rag-test",
	// 	QAIndexName:    "rag-test-qa",
	// 	APIKey:         g.Cfg().MustGet(ctx, "embedding.apiKey").String(),
	// 	BaseURL:        g.Cfg().MustGet(ctx, "embedding.baseURL").String(),
	// 	EmbeddingModel: g.Cfg().MustGet(ctx, "embedding.model").String(),
	// 	ChatModel:      g.Cfg().MustGet(ctx, "chat.model").String(),
	// }
	cfg = &qwen.ChatModelConfig{
		APIKey:  g.Cfg().MustGet(ctx, "qa.apiKey").String(),
		BaseURL: g.Cfg().MustGet(ctx, "qa.baseURL").String(),
		Model:   g.Cfg().MustGet(ctx, "qa.model").String(),
	}
	cm, err := qwen.NewChatModel(ctx, cfg)
	if err != nil {
		return nil, err
	}
	noThinkChatModel = cm
	return cm, nil
}
