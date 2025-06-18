package rag

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/wangle201210/go-rag/server/core"
	"github.com/wangle201210/go-rag/server/core/config"
)

var ragSvr = &core.Rag{}

func init() {
	ctx := gctx.New()
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{g.Cfg().MustGet(ctx, "es.address").String()},
		Username:  g.Cfg().MustGet(ctx, "es.username").String(),
		Password:  g.Cfg().MustGet(ctx, "es.password").String(),
	})
	if err != nil {
		g.Log().Fatalf(ctx, "NewClient of es8 failed, err=%v", err)
		return
	}
	ragSvr, err = core.New(ctx, &config.Config{
		Client:         client,
		IndexName:      g.Cfg().MustGet(ctx, "es.indexName").String(),
		APIKey:         g.Cfg().MustGet(ctx, "embedding.apiKey").String(),
		BaseURL:        g.Cfg().MustGet(ctx, "embedding.baseURL").String(),
		EmbeddingModel: g.Cfg().MustGet(ctx, "embedding.model").String(),
		ChatModel:      g.Cfg().MustGet(ctx, "chat.model").String(),
	})
	if err != nil {
		g.Log().Fatalf(ctx, "New of rag failed, err=%v", err)
		return
	}
}

func GetRagSvr() *core.Rag {
	return ragSvr
}
