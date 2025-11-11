package indexer

import (
	"context"
	"fmt"

	einoqdrant "github.com/cloudwego/eino-ext/components/indexer/qdrant"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/schema"
	"github.com/qdrant/go-client/qdrant"
)

// QdrantIndexerConfig Qdrant indexer 配置
type QdrantIndexerConfig struct {
	Client     *qdrant.Client     // Required: Qdrant client
	Collection string             // Required: Collection name
	VectorDim  int                // Required: Vector dimension
	Distance   qdrant.Distance    // Required: Distance metric
	BatchSize  int                // Optional: Batch size (default: 10)
	Embedding  embedding.Embedder // Required: Embedding component
	IsAsync    bool               // Optional: 是否异步模式（包含 QA 向量）
}

// QdrantIndexer Qdrant indexer 实现（包装 eino-ext 的实现）
type QdrantIndexer struct {
	config       *QdrantIndexerConfig
	einoIndexer  indexer.Indexer // eino-ext 的 indexer
	asyncIndexer indexer.Indexer // 异步模式的 indexer（用于 QA 向量）
}

// NewQdrantIndexer 创建 Qdrant indexer，使用 eino-ext 库
func NewQdrantIndexer(ctx context.Context, config *QdrantIndexerConfig) (indexer.Indexer, error) {
	if config.Client == nil {
		return nil, fmt.Errorf("qdrant client is required")
	}
	if config.Collection == "" {
		return nil, fmt.Errorf("collection name is required")
	}
	if config.Embedding == nil {
		return nil, fmt.Errorf("embedding component is required")
	}
	if config.BatchSize == 0 {
		config.BatchSize = 10
	}
	if config.Distance == 0 {
		config.Distance = qdrant.Distance_Cosine
	}

	// 使用 eino-ext 的 Qdrant indexer
	einoConfig := &einoqdrant.Config{
		Client:     config.Client,
		Collection: config.Collection,
		VectorDim:  config.VectorDim,
		Distance:   config.Distance,
		Embedding:  config.Embedding,
		BatchSize:  config.BatchSize,
	}

	einoIndexer, err := einoqdrant.NewIndexer(ctx, einoConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create eino indexer: %w", err)
	}

	idx := &QdrantIndexer{
		config:      config,
		einoIndexer: einoIndexer,
	}

	// 注意：eino-ext 的 indexer 不支持命名向量
	// 如果需要异步模式（QA 向量），需要使用自定义实现或创建单独的 collection
	// 当前实现：使用自定义逻辑处理命名向量

	return idx, nil
}

// Store 存储文档（实现 Indexer 接口）
// 注意：由于 eino-ext 的 indexer 不支持命名向量，这里使用自定义实现
func (idx *QdrantIndexer) Store(ctx context.Context, docs []*schema.Document, opts ...indexer.Option) ([]string, error) {
	// 使用自定义实现支持命名向量
	return idx.StoreWithNamedVectors(ctx, docs, opts...)
}

// GetType 返回 indexer 类型
func (idx *QdrantIndexer) GetType() string {
	return "qdrant_indexer"
}
