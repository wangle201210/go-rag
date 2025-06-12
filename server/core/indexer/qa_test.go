package indexer

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/wangle201210/go-rag/server/core/common"
	"github.com/wangle201210/go-rag/server/core/config"
)

func TestQA(t *testing.T) {
	ctx := gctx.New()
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	cfg := &config.Config{
		Client:    client,
		IndexName: "rag-test",
		// QAIndexName:    "rag-test-qa",
		APIKey:         g.Cfg().MustGet(ctx, "embedding.apiKey").String(),
		BaseURL:        g.Cfg().MustGet(ctx, "embedding.baseURL").String(),
		EmbeddingModel: g.Cfg().MustGet(ctx, "embedding.model").String(),
		ChatModel:      g.Cfg().MustGet(ctx, "chat.model").String(),
	}
	err = InitQaIndexer(ctx, cfg)
	if err != nil {
		t.Fatal(err)
		return
	}
	ctx = context.WithValue(ctx, common.KnowledgeName, "test")
	err = docQA(ctx, cfg, &schema.Document{
		ID:      "abcd-1234",
		Content: "h1:多语言支持 h2:支持语言 \n\n## 支持语言\n\n| 语言   | 代码       | 状态   |\n\n|------|----------|------|\n\n| 简体中文 | zh-CN    | 完整支持 |\n\n| 繁体中文 | zh-TW/HK | 完整支持 |\n\n| 英语   | en-US    | 完整支持 |\n\n| 法语   | fr-FR    | 完整支持 |\n\n| 日语   | ja-JP    | 完整支持 |\n\n| 韩语   | ko-KR    | 完整支持 |\n\n| 俄语   | ru-RU    | 完整支持 |",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}
