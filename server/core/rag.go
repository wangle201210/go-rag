package core

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/retriever/es8"
	"github.com/cloudwego/eino/components/document"
	er "github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/wangle201210/go-rag/server/core/common"
	"github.com/wangle201210/go-rag/server/core/config"
	"github.com/wangle201210/go-rag/server/core/indexer"
	"github.com/wangle201210/go-rag/server/core/retriever"
)

type Rag struct {
	idxer  compose.Runnable[any, []string]
	rtrvr  compose.Runnable[string, []*schema.Document]
	client *elasticsearch.Client
}

func New(ctx context.Context, conf *config.Config) (*Rag, error) {
	var err error
	// Init eino devops server
	// err := devops.Init(ctx)
	// if err != nil {
	// 	g.Log().Errorf(ctx, "[eino dev] init failed, err=%v", err)
	// 	return nil, err
	// }

	if len(conf.IndexName) == 0 {
		return nil, fmt.Errorf("indexName is empty")
	}
	// 确保es index存在
	err = common.CreateIndexIfNotExists(ctx, conf.Client, conf.IndexName)
	if err != nil {
		return nil, err
	}
	buildIndex, err := indexer.BuildIndexer(ctx, conf)
	if err != nil {
		return nil, err
	}
	buildRetriever, err := retriever.BuildRetriever(ctx, conf)
	if err != nil {
		return nil, err
	}
	return &Rag{
		idxer: buildIndex,
		rtrvr: buildRetriever,
	}, nil
}

type IndexReq struct {
	URI           string // 文档地址，可以是文件路径（pdf，html，md等），也可以是网址
	KnowledgeName string // 知识库名称
}

// Index
// uri:
// ids: 文档id
func (x *Rag) Index(ctx context.Context, req *IndexReq) (ids []string, err error) {
	s := document.Source{
		URI: req.URI,
	}
	ctx = context.WithValue(ctx, common.KnowledgeName, req.KnowledgeName)
	ids, err = x.idxer.Invoke(ctx, s, indexer.WithCallbacks()...)
	if err != nil {
		return
	}
	return
}

type RetrieveReq struct {
	Query         string  // 检索关键词
	TopK          int     // 检索结果数量
	Score         float64 //  分数阀值(0-2, 0 完全相反，1 毫不相干，2 完全相同,一般需要传入一个大于1的数字，如1.5)
	KnowledgeName string  // 知识库名字
}

// Retrieve 检索
func (x *Rag) Retrieve(ctx context.Context, req *RetrieveReq) (msg []*schema.Document, err error) {
	msg, err = x.rtrvr.Invoke(ctx, req.Query,
		compose.WithRetrieverOption(
			er.WithScoreThreshold(req.Score),
			er.WithTopK(req.TopK),
			es8.WithFilters([]types.Query{
				{Match: map[string]types.MatchQuery{common.KnowledgeName: {Query: req.KnowledgeName}}},
			}),
		),
	)
	if err != nil {
		return
	}
	return
}
