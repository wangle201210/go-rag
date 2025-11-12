package common

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/exists"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/gogf/gf/v2/frame/g"
	coretypes "github.com/wangle201210/go-rag/server/core/types"
)

// createIndex create index for example in add_documents.go.
// 已废弃：使用 vector.VectorStore 接口替代
// Deprecated: Use vector.VectorStore.CreateIndex instead
func createIndex(ctx context.Context, client *elasticsearch.Client, indexName string) error {
	_, err := create.NewCreateFunc(client)(indexName).Request(&create.Request{
		Mappings: &types.TypeMapping{
			Properties: map[string]types.Property{
				coretypes.FieldContent:  types.NewTextProperty(),
				coretypes.FieldExtra:    types.NewTextProperty(),
				coretypes.KnowledgeName: types.NewKeywordProperty(),
				coretypes.FieldContentVector: &types.DenseVectorProperty{
					Dims:       Of(1024), // same as embedding dimensions
					Index:      Of(true),
					Similarity: Of("cosine"),
				},
				coretypes.FieldQAContentVector: &types.DenseVectorProperty{
					Dims:       Of(1024), // same as embedding dimensions
					Index:      Of(true),
					Similarity: Of("cosine"),
				},
			},
		},
	}).Do(ctx)

	return err
}

// CreateIndexIfNotExists 已废弃：使用 vector.VectorStore 接口替代
// Deprecated: Use vector.VectorStore.IndexExists and CreateIndex instead
func CreateIndexIfNotExists(ctx context.Context, client *elasticsearch.Client, indexName string) error {
	indexExists, err := exists.NewExistsFunc(client)(indexName).Do(ctx)
	if err != nil {
		return err
	}
	if indexExists {
		return nil
	}
	err = createIndex(ctx, client, indexName)
	return err
}

// DeleteDocument 删除索引中的单个文档
// 已废弃：使用 vector.VectorStore 接口替代
// Deprecated: Use vector.VectorStore.DeleteDocument instead
func DeleteDocument(ctx context.Context, client *elasticsearch.Client, documentID string) error {
	return withRetry(func() error {
		indexName := g.Cfg().MustGet(ctx, "vector.indexName").String()
		res, err := client.Delete(indexName, documentID)
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

// withRetry 包装函数，添加重试机制
func withRetry(operation func() error) error {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 30 * time.Second

	return backoff.Retry(operation, b)
}
