package rag

import (
	"context"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/wangle201210/go-rag"
	"github.com/wangle201210/go-rag/config"
)

var ragSvr = &rag.Rag{}

func init() {
	ctx := context.Background()
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{g.Cfg().MustGet(ctx, "es.address").String()},
	})
	if err != nil {
		log.Printf("NewClient of es8 failed, err=%v", err)
		return
	}
	ragSvr, err = rag.New(ctx, &config.Config{
		Client:    client,
		IndexName: g.Cfg().MustGet(ctx, "es.indexName").String(),
		APIKey:    g.Cfg().MustGet(ctx, "embedding.apiKey").String(),
		BaseURL:   g.Cfg().MustGet(ctx, "embedding.baseURL").String(),
		Model:     g.Cfg().MustGet(ctx, "embedding.model").String(),
	})
	if err != nil {
		log.Printf("New of rag failed, err=%v", err)
		return
	}
}

func GetRagSvr() *rag.Rag {
	return ragSvr
}
