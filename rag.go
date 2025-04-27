package rag

import (
	"context"

	"github.com/cloudwego/eino/components/document"
	er "github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/wangle201210/go-rag/indexer"
	"github.com/wangle201210/go-rag/retriever"
)

// 不建议直接使用下面的方法，这样没调用一次都会new一次对象
// 建议自己调用 indexer.BuildIndexer 或者 retriever.BuildRetriever，这样可以复用

// Index
// uri: 文档地址，可以是文件路径（pdf，html，md等），也可以是网址
// ids: 文档id
func Index(uri string) (ids []string, err error) {
	buildIndex, err := indexer.BuildIndexer(context.Background())
	if err != nil {
		return
	}
	s := document.Source{
		URI: uri,
	}
	ids, err = buildIndex.Invoke(context.Background(), s)
	if err != nil {
		return
	}
	return
}

// Retrieve
// input: 检索关键词
// score: 0-2, 0 完全相反，1 毫不相干，2 完全相同
func Retrieve(input string, score float64) (msg []*schema.Document, err error) {
	r, err := retriever.BuildRetriever(context.Background())
	if err != nil {
		return
	}
	msg, err = r.Invoke(context.Background(), input,
		compose.WithRetrieverOption(
			er.WithScoreThreshold(score),
			er.WithTopK(5),
		),
	)
	if err != nil {
		return
	}
	return
}
