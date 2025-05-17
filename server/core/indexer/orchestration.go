package indexer

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/wangle201210/go-rag/server/core/config"
)

func BuildIndexer(ctx context.Context, conf *config.Config) (r compose.Runnable[any, []string], err error) {
	const (
		Loader1              = "Loader1"
		Indexer2             = "Indexer2"
		DocumentTransformer3 = "DocumentTransformer3"
	)
	g := compose.NewGraph[any, []string]()
	loader1KeyOfLoader, err := newLoader(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddLoaderNode(Loader1, loader1KeyOfLoader)
	indexer2KeyOfIndexer, err := newIndexer(ctx, conf)
	if err != nil {
		return nil, err
	}
	_ = g.AddIndexerNode(Indexer2, indexer2KeyOfIndexer)
	documentTransformer2KeyOfDocumentTransformer, err := newDocumentTransformer(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddDocumentTransformerNode(DocumentTransformer3, documentTransformer2KeyOfDocumentTransformer)
	_ = g.AddEdge(compose.START, Loader1)
	_ = g.AddEdge(Loader1, DocumentTransformer3)
	_ = g.AddEdge(DocumentTransformer3, Indexer2)
	_ = g.AddEdge(Indexer2, compose.END)
	r, err = g.Compile(ctx, compose.WithGraphName("indexer"))
	if err != nil {
		return nil, err
	}
	return r, err
}
