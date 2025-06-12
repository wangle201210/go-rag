package config

import (
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/elastic/go-elasticsearch/v8"
)

type Config struct {
	Client    *elasticsearch.Client
	IndexName string // es index name
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
		Client:    x.Client,
		IndexName: x.IndexName,
		// embedding 时使用
		APIKey:         x.APIKey,
		BaseURL:        x.BaseURL,
		EmbeddingModel: x.EmbeddingModel,
		ChatModel:      x.ChatModel,
	}
}
