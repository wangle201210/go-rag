package rag

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/wangle201210/go-rag/config"
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
	ids, err := ragSvr.Index("./test_file/readme.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range ids {
		t.Log(id)
	}
	ragSvr.Index("./test_file/readme2.md")
	ragSvr.Index("./test_file/readme.html")
	ragSvr.Index("./test_file/test.pdf")
	ragSvr.Index("https://deepchat.thinkinai.xyz/docs/guide/advanced-features/shortcuts.html")
}

func TestRetriever(t *testing.T) {
	msg, err := ragSvr.Retrieve("这里有很多内容", 1.5, 5)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range msg {
		t.Logf("content: %v, score: %v", m.Content, m.Score())
	}

	msg, err = ragSvr.Retrieve("代码解析", 1.5, 5)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range msg {
		t.Logf(" content: %v, score: %v", m.Content, m.Score())
	}
}
