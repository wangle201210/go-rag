package retriever

import (
	"context"

	"github.com/cloudwego/eino-ext/components/document/transformer/reranker/score"
	"github.com/cloudwego/eino/components/document"
)

// newDocumentTransformer component initialization function of node 'DocumentTransformer1' in graph 'retriever'
func newDocumentTransformer(ctx context.Context) (tfr document.Transformer, err error) {
	config := &score.Config{}
	tfr, err = score.NewReranker(ctx, config)
	if err != nil {
		return nil, err
	}
	return tfr, nil
}
