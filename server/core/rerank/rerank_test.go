package rerank

import (
	"os"
	"testing"

	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/os/gctx"
)

func TestRerank(t *testing.T) {
	// 该测试会真实调用 rerank 在线接口，未配置 RERANK_API_KEY 时跳过，
	// 保证 go test ./... 可以离线通过
	apiKey := os.Getenv("RERANK_API_KEY")
	if apiKey == "" {
		t.Skip("未设置环境变量 RERANK_API_KEY，跳过 rerank 在线测试")
	}
	rerankCfg = &Conf{
		apiKey:          apiKey,
		Model:           "BAAI/bge-reranker-v2-m3",
		ReturnDocuments: false,
		MaxChunksPerDoc: 1024,
		OverlapTokens:   80,
		url:             "https://api.siliconflow.cn/v1/rerank",
	}
	ctx := gctx.New()
	docs := []*schema.Document{
		{Content: "banana"},
		{Content: "fruit"},
		{Content: "apple"},
		{Content: "vegetable"},
	}
	output, err := NewRerank(ctx, "水果", docs, 2)
	if err != nil {
		t.Fatal(err)
	}
	for _, doc := range output {
		t.Logf("content: %v, score: %v", doc.Content, doc.Score())
	}
}
