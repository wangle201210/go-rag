package common

import (
	"context"
	"os"
	"strings"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/wangle201210/go-rag/server/core/config"
)

func NewEmbedding(ctx context.Context, conf *config.Config) (eb embedding.Embedder, err error) {
	econf := &openai.EmbeddingConfig{
		APIKey:  conf.APIKey,
		Model:   conf.EmbeddingModel,
		Timeout: 0,
		BaseURL: conf.BaseURL,
	}
	if econf.APIKey == "" {
		econf.APIKey = os.Getenv("OPENAI_API_KEY")
	}
	if econf.BaseURL == "" {
		econf.BaseURL = os.Getenv("OPENAI_BASE_URL")
	}
	if econf.Model == "" {
		econf.Model = "text-embedding-3-large"
	}
	// Only request dimensions for models that support it. BGE-M3 already
	// returns 1024 dimensions and SiliconFlow rejects this optional parameter.
	if strings.HasPrefix(econf.Model, "text-embedding-3-") || strings.HasPrefix(econf.Model, "Qwen/Qwen3-Embedding-") {
		econf.Dimensions = Of(1024)
	}
	eb, err = openai.NewEmbedder(ctx, econf)
	if err != nil {
		return nil, err
	}
	return eb, nil
}
