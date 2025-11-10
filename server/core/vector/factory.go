package vector

import (
	"fmt"
)

// NewVectorStore 创建向量存储实例
func NewVectorStore(cfg *Config) (VectorStore, error) {
	switch cfg.Type {
	case "es", "elasticsearch":
		if cfg.ES == nil {
			return nil, fmt.Errorf("es config is required when type is es")
		}
		return NewESVectorStore(cfg.ES)
	case "qdrant":
		if cfg.Qdrant == nil {
			return nil, fmt.Errorf("qdrant config is required when type is qdrant")
		}
		return NewQdrantVectorStore(cfg.Qdrant)
	default:
		return nil, fmt.Errorf("unsupported vector store type: %s", cfg.Type)
	}
}
