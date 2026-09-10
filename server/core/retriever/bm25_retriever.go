package retriever

import (
	"context"

	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/gogf/gf/v2/errors/gerror"
	coretypes "github.com/wangle201210/go-rag/server/core/types"
)

// Bm25Retrieve 使用 ES BM25 对 content 字段做全文检索。
// 仅在使用 ES 向量存储时调用，复用同一索引中的文档。
func Bm25Retrieve(ctx context.Context, client *elasticsearch.Client, indexName, knowledgeName, query string, excludeIDs []string, topK int) ([]*schema.Document, error) {
	if client == nil {
		return nil, gerror.New("bm25 requires an ES client")
	}
	if topK <= 0 {
		return nil, gerror.New("bm25 topK must be positive")
	}
	boolQuery := &types.BoolQuery{
		Must: []types.Query{{Match: map[string]types.MatchQuery{
			coretypes.FieldContent: {Query: query},
		}}},
		Filter: []types.Query{{Term: map[string]types.TermQuery{
			coretypes.KnowledgeName: {Value: knowledgeName},
		}}},
	}

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
		return nil, gerror.Wrap(err, "bm25 search failed")
	}

	var docs []*schema.Document
	for _, hit := range resp.Hits.Hits {
		doc, err := EsHit2Document(ctx, hit)
		if err != nil {
			return nil, gerror.Wrap(err, "parse bm25 search result")
		}
		docs = append(docs, doc)
	}

	return docs, nil
}
