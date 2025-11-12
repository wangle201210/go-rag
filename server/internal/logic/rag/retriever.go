package rag

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/qdrant/go-client/qdrant"
	"github.com/wangle201210/go-rag/server/core"
	"github.com/wangle201210/go-rag/server/core/config"
	"github.com/wangle201210/go-rag/server/core/vector"
)

var ragSvr = &core.Rag{}

func init() {
	ctx := gctx.New()

	// 读取向量存储配置
	vectorType := g.Cfg().MustGet(ctx, "vector.type").String()
	indexName := g.Cfg().MustGet(ctx, "vector.indexName").String()

	// 创建向量存储配置
	vectorCfg := &vector.Config{
		Type:      vectorType,
		IndexName: indexName,
	}

	// 根据类型配置
	if vectorType == "es" || vectorType == "elasticsearch" {
		vectorCfg.ES = &vector.ESConfig{
			Address:  g.Cfg().MustGet(ctx, "vector.es.address").String(),
			Username: g.Cfg().MustGet(ctx, "vector.es.username").String(),
			Password: g.Cfg().MustGet(ctx, "vector.es.password").String(),
		}
	} else if vectorType == "qdrant" {
		vectorCfg.Qdrant = &vector.QdrantConfig{
			Address: g.Cfg().MustGet(ctx, "vector.qdrant.address").String(),
			Port:    g.Cfg().MustGet(ctx, "vector.qdrant.port").Int(),
			APIKey:  g.Cfg().MustGet(ctx, "vector.qdrant.apiKey").String(),
		}
	}

	// 创建向量存储实例
	vectorStore, err := vector.NewVectorStore(vectorCfg)
	if err != nil {
		g.Log().Fatalf(ctx, "NewVectorStore failed, err=%v", err)
		return
	}

	// 根据类型获取对应的客户端
	var client *elasticsearch.Client
	var qdrantClient *qdrant.Client

	if esStore, ok := vectorStore.(*vector.ESVectorStore); ok {
		client = esStore.GetClient()
	} else if qdrantStore, ok := vectorStore.(*vector.QdrantVectorStore); ok {
		qdrantClient = qdrantStore.GetClient()
	}

	ragSvr, err = core.New(ctx, &config.Config{
		Client:         client,
		QdrantClient:   qdrantClient,
		IndexName:      indexName,
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
