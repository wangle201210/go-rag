package rerank

import (
	"testing"

	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/os/gctx"
)

func TestRerank(t *testing.T) {
	rerankCfg = &Conf{
		apiKey:          "sk-***",
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
	output, err := Rerank(ctx, "水果", docs, 2)
	if err != nil {
		t.Fatal(err)
	}
	for _, doc := range output {
		t.Logf("content: %v, score: %v", doc.Content, doc.Score())
	}
}
