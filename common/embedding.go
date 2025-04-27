package common

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino/components/embedding"
)

func NewEmbedding(ctx context.Context) (eb embedding.Embedder, err error) {
	config := &openai.EmbeddingConfig{
		APIKey:     os.Getenv("OPENAI_API_KEY"),
		Model:      "text-embedding-3-large",
		Dimensions: Of(1024),
		Timeout:    0,
		BaseURL:    os.Getenv("OPENAI_BASE_URL"),
	}
	eb, err = openai.NewEmbedder(ctx, config)
	if err != nil {
		return nil, err
	}
	return eb, nil
}
