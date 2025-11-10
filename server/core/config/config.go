package config

import (
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/wangle201210/go-rag/server/core/vector"
)

type Config struct {
	Client      *elasticsearch.Client // 保留用于兼容
	VectorStore vector.VectorStore    // 新的向量存储接口
	IndexName   string                // index name
	// embedding 时使用
	APIKey         string
	BaseURL        string
	EmbeddingModel string
	ChatModel      string
}

func (x *Config) GetChatModelConfig() *openai.ChatModelConfig {
	if x == nil {
		return nil
	}
	return &openai.ChatModelConfig{
		APIKey:  x.APIKey,
		BaseURL: x.BaseURL,
		Model:   x.ChatModel,
	}
}

func (x *Config) Copy() *Config {
	return &Config{
		Client:      x.Client,
		VectorStore: x.VectorStore,
		IndexName:   x.IndexName,
		// embedding 时使用
		APIKey:         x.APIKey,
		BaseURL:        x.BaseURL,
		EmbeddingModel: x.EmbeddingModel,
		ChatModel:      x.ChatModel,
	}
}
