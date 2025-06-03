package core

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/wangle201210/go-rag/server/core/config"
)

var ragSvr = &Rag{}

func init() {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		log.Printf("NewClient of es8 failed, err=%v", err)
		return
	}
	ragSvr, err = New(context.Background(), &config.Config{
		Client:    client,
		IndexName: "rag-test",
		APIKey:    os.Getenv("OPENAI_API_KEY"),
		BaseURL:   os.Getenv("OPENAI_BASE_URL"),
		Model:     "text-embedding-3-large",
	})
	if err != nil {
		log.Printf("New of rag failed, err=%v", err)
		return
	}
}
func TestIndex(t *testing.T) {
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
	ctx := context.Background()
	req := &RetrieveReq{
		Query:         "这里有很多内容",
		TopK:          5,
		Score:         1.2,
		KnowledgeName: "wanna",
	}
	msg, err := ragSvr.Retrieve(ctx, req)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range msg {
		t.Logf("content: %v, score: %v", m.Content, m.Score())
	}
}
