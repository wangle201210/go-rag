package retriever

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	coretypes "github.com/wangle201210/go-rag/server/core/types"
)

// Bm25Retrieve 使用 ES BM25 对 content 字段做全文检索。
// 仅在 conf.ESClient != nil 时调用。
func Bm25Retrieve(ctx context.Context, client *elasticsearch.Client, indexName, knowledgeName, query string, excludeIDs []string, topK int) ([]*schema.Document, error) {
	must := []types.Query{
		{Match: map[string]types.MatchQuery{
			coretypes.FieldContent: {Query: query},
		}},
		{Bool: &types.BoolQuery{
			Must: []types.Query{
				{Match: map[string]types.MatchQuery{
					coretypes.KnowledgeName: {Query: knowledgeName},
				}},
			},
		}},
	}

	boolQuery := &types.BoolQuery{Must: must}

	if len(excludeIDs) > 0 {
		boolQuery.MustNot = []types.Query{
			{Terms: &types.TermsQuery{
				TermsQuery: map[string]types.TermsQueryField{
					"_id": excludeIDs,
				},
			}},
		}
	}

	sreq := search.NewRequest()
	sreq.Query = &types.Query{Bool: boolQuery}
	sreq.Size = &topK

	resp, err := search.NewSearchFunc(client)().
		Index(indexName).
		Request(sreq).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("bm25 search failed: %w", err)
	}

	var docs []*schema.Document
	for _, hit := range resp.Hits.Hits {
		doc, err := EsHit2Document(ctx, hit)
		if err != nil {
			continue
		}
		docs = append(docs, doc)
	}

	return docs, nil
}
