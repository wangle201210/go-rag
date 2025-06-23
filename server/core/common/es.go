package common

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gogf/gf/v2/frame/g"
)

// createIndex create index for example in add_documents.go.
func createIndex(ctx context.Context, client *elasticsearch.Client, indexName string) error {
	settings := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   3,
			"number_of_replicas": 1,
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"default": map[string]interface{}{
						"type": "standard",
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				FieldContent:  map[string]interface{}{"type": "text"},
				FieldExtra:    map[string]interface{}{"type": "text"},
				KnowledgeName: map[string]interface{}{"type": "keyword"},
				FieldContentVector: map[string]interface{}{
					"type":       "dense_vector",
					"dims":       1024,
					"index":      true,
					"similarity": "cosine",
				},
				FieldQAContentVector: map[string]interface{}{
					"type":       "dense_vector",
					"dims":       1024,
					"index":      true,
					"similarity": "cosine",
				},
			},
		},
	}

	res, err := client.Indices.Create(
		indexName,
		client.Indices.Create.WithContext(ctx),
		client.Indices.Create.WithBody(strings.NewReader(toJSON(settings))),
	)
	if err != nil {
		return fmt.Errorf("create index failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("create index failed: %s", res.String())
	}

	return nil
}

// withRetry 包装函数，添加重试机制
func withRetry(operation func() error) error {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 30 * time.Second

	return backoff.Retry(operation, b)
}

func CreateIndexIfNotExists(ctx context.Context, client *elasticsearch.Client, indexName string) error {
	return withRetry(func() error {
		exists, err := client.Indices.Exists([]string{indexName})
		if err != nil {
			return fmt.Errorf("check index exists failed: %w", err)
		}
		if exists.StatusCode == 200 {
			return nil
		}
		if err = createIndex(ctx, client, indexName); err != nil {
			return fmt.Errorf("create index failed: %w", err)
		}
		return nil
	})
}

// DeleteDocument 删除索引中的单个文档
func DeleteDocument(ctx context.Context, client *elasticsearch.Client, documentID string) error {
	return withRetry(func() error {
		indexName := g.Cfg().MustGet(ctx, "es.indexName").String()
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

// toJSON 将对象转换为JSON字符串
func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
