package vector

import (
	"context"
	"testing"
)

func TestQdrantVectorStore_CreateIndex(t *testing.T) {
	// 这是一个示例测试，需要实际的 Qdrant 服务器才能运行
	t.Skip("需要实际的 Qdrant 服务器")

	cfg := &QdrantConfig{
		Address: "localhost:6334",
		APIKey:  "",
	}

	store, err := NewQdrantVectorStore(cfg)
	if err != nil {
		t.Fatalf("创建 Qdrant 存储失败: %v", err)
	}
	defer store.Close()

	ctx := context.Background()
	indexName := "test_index"

	// 测试创建索引
	err = store.CreateIndex(ctx, indexName)
	if err != nil {
		t.Fatalf("创建索引失败: %v", err)
	}

	// 测试索引是否存在
	exists, err := store.IndexExists(ctx, indexName)
	if err != nil {
		t.Fatalf("检查索引失败: %v", err)
	}

	if !exists {
		t.Errorf("索引应该存在")
	}
}
