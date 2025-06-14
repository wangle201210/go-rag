package retriever

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/wangle201210/go-rag/server/core/config"
)

func BuildRetriever(ctx context.Context, conf *config.Config) (r compose.Runnable[string, []*schema.Document], err error) {
	const (
		Retriever1 = "Retriever"
	)
	g := compose.NewGraph[string, []*schema.Document]()
	retriever1KeyOfRetriever, err := newRetriever(ctx, conf)
	if err != nil {
		return nil, err
	}
	_ = g.AddRetrieverNode(Retriever1, retriever1KeyOfRetriever)
	_ = g.AddEdge(compose.START, Retriever1)
	_ = g.AddEdge(Retriever1, compose.END)
	r, err = g.Compile(ctx, compose.WithGraphName("retriever"))
	if err != nil {
		return nil, err
	}
	return r, err
}
