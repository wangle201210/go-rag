package retriever

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
	"github.com/qdrant/go-client/qdrant"
	"github.com/wangle201210/go-rag/server/core/vector"
)

// QdrantRetrieverConfig Qdrant retriever 配置
type QdrantRetrieverConfig struct {
	VectorStore *vector.QdrantVectorStore
	IndexName   string
	VectorField string
	Embedding   embedding.Embedder
	TopK        int
}

// QdrantRetriever Qdrant retriever 实现
type QdrantRetriever struct {
	config *QdrantRetrieverConfig
}

// NewQdrantRetriever 创建 Qdrant retriever
func NewQdrantRetriever(ctx context.Context, config *QdrantRetrieverConfig) (retriever.Retriever, error) {
	if config.TopK == 0 {
		config.TopK = 10
	}
	return &QdrantRetriever{
		config: config,
	}, nil
}

// Retrieve 检索文档
func (r *QdrantRetriever) Retrieve(ctx context.Context, query string, opts ...retriever.Option) ([]*schema.Document, error) {
	// 解析选项
	options := &retriever.Options{}
	retriever.GetCommonOptions(options, opts...)

	// 获取查询向量
	embedResp, err := r.config.Embedding.EmbedStrings(ctx, []string{query})
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	if len(embedResp) == 0 || len(embedResp[0]) == 0 {
		return nil, fmt.Errorf("empty embedding result")
	}

	queryVector := embedResp[0]

	// 转换为 float32
	queryVec32 := make([]float32, len(queryVector))
	for i, v := range queryVector {
		queryVec32[i] = float32(v)
	}

	// 构建 Qdrant 查询
	topK := r.config.TopK
	if options.TopK != nil && *options.TopK > 0 {
		topK = *options.TopK
	}

	qdrantQuery := &vector.QdrantSearchQuery{
		CollectionName: r.config.IndexName,
		Query: &qdrant.Query{
			Variant: &qdrant.Query_Nearest{
				Nearest: &qdrant.VectorInput{
					Variant: &qdrant.VectorInput_Dense{
						Dense: &qdrant.DenseVector{
							Data: queryVec32,
						},
					},
				},
			},
		},
	}

	// 执行搜索
	searchReq := &vector.SearchRequest{
		IndexName: r.config.IndexName,
		Query:     qdrantQuery,
		Size:      topK,
	}

	resp, err := r.config.VectorStore.SearchDocuments(ctx, searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to search documents: %w", err)
	}

	return resp.Documents, nil
}

// GetType 返回 retriever 类型
func (r *QdrantRetriever) GetType() string {
	return "qdrant_retriever"
}
