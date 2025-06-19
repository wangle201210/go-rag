package indexer

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/wangle201210/go-rag/server/core/config"
)

func BuildIndexerAsync(ctx context.Context, conf *config.Config) (r compose.Runnable[[]*schema.Document, []string], err error) {
	const (
		Indexer = "Indexer"
		QA      = "QA"
	)

	g := compose.NewGraph[[]*schema.Document, []string]()
	indexer2KeyOfIndexer, err := newAsyncIndexer(ctx, conf)
	if err != nil {
		return nil, err
	}
	_ = g.AddIndexerNode(Indexer, indexer2KeyOfIndexer)
	_ = g.AddLambdaNode(QA, compose.InvokableLambda(qa))
	_ = g.AddEdge(compose.START, QA)
	_ = g.AddEdge(QA, Indexer)
	_ = g.AddEdge(Indexer, compose.END)
	r, err = g.Compile(ctx, compose.WithGraphName("indexer_async"))
	if err != nil {
		return nil, err
	}
	return r, err
}
