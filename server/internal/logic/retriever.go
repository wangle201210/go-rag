package logic

import (
	"context"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/wangle201210/go-rag"
	"github.com/wangle201210/go-rag/config"
)

var ragSvr = &rag.Rag{}

func init() {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		log.Printf("NewClient of es8 failed, err=%v", err)
		return
	}
	ragSvr, err = rag.New(context.Background(), &config.Config{
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

func GetRagSvr() *rag.Rag {
	return ragSvr
}
