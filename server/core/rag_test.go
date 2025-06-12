package core

import (
	"context"
	"log"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/wangle201210/go-rag/server/core/config"
)

var ragSvr = &Rag{}
var cfg = &config.Config{}

func _init() {
	ctx := context.Background()
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		log.Printf("NewClient of es8 failed, err=%v", err)
		return
	}
	cfg = &config.Config{
		Client:         client,
		IndexName:      "rag-test1",
		APIKey:         g.Cfg().MustGet(ctx, "embedding.apiKey").String(),
		BaseURL:        g.Cfg().MustGet(ctx, "embedding.baseURL").String(),
		EmbeddingModel: g.Cfg().MustGet(ctx, "embedding.model").String(),
		ChatModel:      g.Cfg().MustGet(ctx, "chat.model").String(),
	}
	ragSvr, err = New(context.Background(), cfg)
	if err != nil {
		log.Printf("New of rag failed, err=%v", err)
		return
	}
}
func TestIndex(t *testing.T) {
	_init()
	ctx := context.Background()
	uriList := []string{
		"./test_file/readme.md",
		// "./test_file/readme2.md",
		// "./test_file/readme.html",
		// "./test_file/test.pdf",
		// "https://deepchat.thinkinai.xyz/docs/guide/advanced-features/shortcuts.html",
	}
	for _, s := range uriList {
		req := &IndexReq{
			URI:           s,
			KnowledgeName: "wanna",
		}
		ids, err := ragSvr.Index(ctx, req)
		if err != nil {
			t.Fatal(err)
		}
		for _, id := range ids {
			t.Log(id)
		}
	}
}

func TestRetriever(t *testing.T) {
	_init()
	ctx := context.Background()
	req := &RetrieveReq{
		Query:         "dify知识库配置",
		TopK:          5,
		Score:         1.2,
		KnowledgeName: "deepchat使用文档",
	}
	msg, err := ragSvr.Retrieve(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range msg {
		t.Logf("content: %v, score: %v", m.Content, m.Score())
	}
}

func TestRag_GetKnowledgeList(t *testing.T) {
	ctx := context.Background()
	list, err := ragSvr.GetKnowledgeBaseList(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("list: %v", list)
}
