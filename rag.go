package rag

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/document"
	er "github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/wangle201210/go-rag/common"
	"github.com/wangle201210/go-rag/config"
	"github.com/wangle201210/go-rag/indexer"
	"github.com/wangle201210/go-rag/retriever"
)

type Rag struct {
	idxer  compose.Runnable[any, []string]
	rtrvr  compose.Runnable[string, []*schema.Document]
	client *elasticsearch.Client
}

func New(ctx context.Context, conf *config.Config) (*Rag, error) {
	if len(conf.IndexName) == 0 {
		return nil, fmt.Errorf("indexName is empty")
	}
	// 确保es index存在
	err := common.CreateIndexIfNotExists(ctx, conf.Client, conf.IndexName)
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

// Index
// uri: 文档地址，可以是文件路径（pdf，html，md等），也可以是网址
// ids: 文档id
func (x *Rag) Index(uri string) (ids []string, err error) {
	s := document.Source{
		URI: uri,
	}
	ids, err = x.idxer.Invoke(context.Background(), s)
	if err != nil {
		return
	}
	return
}

// Retrieve
// input: 检索关键词
// score: 分数阀值(0-2, 0 完全相反，1 毫不相干，2 完全相同,一般需要传入一个大于1的数字，如1.5)
// topK: 检索结果数量
func (x *Rag) Retrieve(input string, score float64, topK int) (msg []*schema.Document, err error) {
	msg, err = x.rtrvr.Invoke(context.Background(), input,
		compose.WithRetrieverOption(
			er.WithScoreThreshold(score),
			er.WithTopK(topK),
		),
	)
	if err != nil {
		return
	}
	return
}
