package vector

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cenkalti/backoff/v4"
	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/exists"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/gogf/gf/v2/frame/g"
	coretypes "github.com/wangle201210/go-rag/server/core/types"
)

// ESVectorStore Elasticsearch 向量存储实现
type ESVectorStore struct {
	client *elasticsearch.Client
}

// NewESVectorStore 创建 ES 向量存储实例
func NewESVectorStore(cfg *ESConfig) (*ESVectorStore, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.Address},
		Username:  cfg.Username,
		Password:  cfg.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create es client: %w", err)
	}

	return &ESVectorStore{
		client: client,
	}, nil
}

// GetClient 获取 ES 客户端（用于兼容现有代码）
func (s *ESVectorStore) GetClient() *elasticsearch.Client {
	return s.client
}

// CreateIndex 创建索引
func (s *ESVectorStore) CreateIndex(ctx context.Context, indexName string) error {
	_, err := create.NewCreateFunc(s.client)(indexName).Request(&create.Request{
		Mappings: &types.TypeMapping{
			Properties: map[string]types.Property{
				coretypes.FieldContent:  types.NewTextProperty(),
				coretypes.FieldExtra:    types.NewTextProperty(),
				coretypes.KnowledgeName: types.NewKeywordProperty(),
				coretypes.FieldContentVector: &types.DenseVectorProperty{
					Dims:       of(1024),
					Index:      of(true),
					Similarity: of("cosine"),
				},
				coretypes.FieldQAContentVector: &types.DenseVectorProperty{
					Dims:       of(1024),
					Index:      of(true),
					Similarity: of("cosine"),
				},
			},
		},
	}).Do(ctx)

	return err
}

// IndexExists 检查索引是否存在
func (s *ESVectorStore) IndexExists(ctx context.Context, indexName string) (bool, error) {
	return exists.NewExistsFunc(s.client)(indexName).Do(ctx)
}

// DeleteDocument 删除文档
func (s *ESVectorStore) DeleteDocument(ctx context.Context, indexName, documentID string) error {
	return s.withRetry(func() error {
		res, err := s.client.Delete(indexName, documentID)
		if err != nil {
			return fmt.Errorf("delete document failed: %w", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			return fmt.Errorf("delete document failed: %s", res.String())
		}

		return nil
	})
}

// GetKnowledgeBaseList 获取知识库列表
func (s *ESVectorStore) GetKnowledgeBaseList(ctx context.Context, indexName string) ([]string, error) {
	names := "distinct_knowledge_names"
	query := search.NewRequest()
	query.Size = of(0)
	query.Aggregations = map[string]types.Aggregations{
		names: {
			Terms: &types.TermsAggregation{
				Field: of(coretypes.KnowledgeName),
				Size:  of(10000),
			},
		},
	}

	res, err := search.NewSearchFunc(s.client)().
		Index(indexName).
		Request(query).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	if res.Aggregations == nil {
		g.Log().Infof(ctx, "No aggregations found")
		return nil, nil
	}

	termsAgg, ok := res.Aggregations[names].(*types.StringTermsAggregate)
	if !ok || termsAgg == nil {
		return nil, errors.New("failed to parse terms aggregation")
	}

	var list []string
	for _, bucket := range termsAgg.Buckets.([]types.StringTermsBucket) {
		list = append(list, bucket.Key.(string))
	}

	return list, nil
}

// SearchDocuments 搜索文档
func (s *ESVectorStore) SearchDocuments(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	sreq := search.NewRequest()
	sreq.Query = req.Query.(*types.Query)
	sreq.Size = of(req.Size)

	resp, err := search.NewSearchFunc(s.client)().
		Index(req.IndexName).
		Request(sreq).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	var docs []*schema.Document
	for _, hit := range resp.Hits.Hits {
		doc, err := s.esHit2Document(ctx, hit)
		if err != nil {
			g.Log().Errorf(ctx, "esHit2Document failed, err=%v", err)
			continue
		}
		docs = append(docs, doc)
	}

	return &SearchResponse{
		Documents: docs,
		Total:     resp.Hits.Total.Value,
	}, nil
}

// esHit2Document 将 ES Hit 转换为 Document
func (s *ESVectorStore) esHit2Document(ctx context.Context, hit types.Hit) (*schema.Document, error) {
	doc := &schema.Document{
		ID:       *hit.Id_,
		MetaData: map[string]any{},
	}

	var src map[string]any
	if err := sonic.Unmarshal(hit.Source_, &src); err != nil {
		return nil, err
	}

	for field, val := range src {
		switch field {
		case coretypes.FieldContent:
			doc.Content = val.(string)
		case coretypes.FieldContentVector:
			var v []float64
			for _, item := range val.([]interface{}) {
				v = append(v, item.(float64))
			}
			doc.WithDenseVector(v)
		case coretypes.FieldQAContentVector, coretypes.FieldQAContent:
			// 这两个字段都不返回
		case coretypes.FieldExtra:
			if val == nil {
				continue
			}
			doc.MetaData[coretypes.FieldExtra] = val.(string)
		case coretypes.KnowledgeName:
			doc.MetaData[coretypes.KnowledgeName] = val.(string)
		default:
			return nil, fmt.Errorf("unexpected field=%s, val=%v", field, val)
		}
	}

	if hit.Score_ != nil {
		doc.WithScore(float64(*hit.Score_))
	}

	return doc, nil
}

// Close 关闭连接
func (s *ESVectorStore) Close() error {
	// ES client 不需要显式关闭
	return nil
}

// withRetry 包装函数，添加重试机制
func (s *ESVectorStore) withRetry(operation func() error) error {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 30 * time.Second
	return backoff.Retry(operation, b)
}

// of 辅助函数，返回指针
func of[T any](v T) *T {
	return &v
}
