package vector

import (
	"context"

	"github.com/cloudwego/eino/schema"
)

// VectorStore 向量存储接口
type VectorStore interface {
	// CreateIndex 创建索引
	CreateIndex(ctx context.Context, indexName string) error

	// IndexExists 检查索引是否存在
	IndexExists(ctx context.Context, indexName string) (bool, error)

	// DeleteDocument 删除文档
	DeleteDocument(ctx context.Context, indexName, documentID string) error

	// GetKnowledgeBaseList 获取知识库列表
	GetKnowledgeBaseList(ctx context.Context, indexName string) ([]string, error)

	// SearchDocuments 搜索文档
	SearchDocuments(ctx context.Context, req *SearchRequest) (*SearchResponse, error)

	// Close 关闭连接
	Close() error
}

// SearchRequest 搜索请求
type SearchRequest struct {
	IndexName     string
	Query         interface{} // 查询条件，不同实现可能不同
	Size          int
	KnowledgeName string
	DocIDs        []string
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Documents []*schema.Document
	Total     int64
}

// Config 向量存储配置
type Config struct {
	Type      string        // 类型：es 或 qdrant
	IndexName string        // 索引名称
	ES        *ESConfig     // ES 配置
	Qdrant    *QdrantConfig // Qdrant 配置
}

// ESConfig Elasticsearch 配置
type ESConfig struct {
	Address  string
	Username string
	Password string
}

// QdrantConfig Qdrant 配置
type QdrantConfig struct {
	Address string
	Port    int
	APIKey  string
}
