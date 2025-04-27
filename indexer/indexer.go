package indexer

import (
	"context"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/components/indexer/es8"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/wangle201210/go-rag/common"
)

// newIndexer component initialization function of node 'Indexer2' in graph 'rag'
func newIndexer(ctx context.Context) (idr indexer.Indexer, err error) {
	client, err := common.GetESClient(ctx)
	if err != nil {
		return nil, err
	}
	config := &es8.IndexerConfig{
		Client:    client,
		Index:     common.IndexName,
		BatchSize: 10,
		DocumentToFields: func(ctx context.Context, doc *schema.Document) (field2Value map[string]es8.FieldValue, err error) {
			if doc.ID == "" {
				doc.ID = uuid.New().String()
			}
			if doc.MetaData != nil {
				marshal, _ := sonic.Marshal(doc.MetaData)
				doc.MetaData[common.DocExtra] = string(marshal)
			}
			return map[string]es8.FieldValue{
				common.FieldContent: {
					Value:    doc.Content,
					EmbedKey: common.FieldContentVector, // vectorize doc content and save vector to field "content_vector"
				},
				common.FieldExtra: {
					Value: doc.MetaData[common.DocExtra],
				},
			}, nil
		},
	}
	embeddingIns11, err := common.NewEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	config.Embedding = embeddingIns11
	idr, err = es8.NewIndexer(ctx, config)
	if err != nil {
		return nil, err
	}
	return idr, nil
}
