package indexer

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/components/indexer/es8"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/schema"
	"github.com/wangle201210/go-rag/server/core/common"
	"github.com/wangle201210/go-rag/server/core/config"
	coretypes "github.com/wangle201210/go-rag/server/core/types"
)

// newAsyncIndexer component initialization function of node 'Indexer2' in graph 'rag'
func newAsyncIndexer(ctx context.Context, conf *config.Config) (idr indexer.Indexer, err error) {
	embeddingIns11, err := common.NewEmbedding(ctx, conf)
	if err != nil {
		return nil, err
	}

	// 根据向量存储类型创建不同的 indexer
	if conf.Client != nil {
		// ES indexer
		indexerConfig := &es8.IndexerConfig{
			Client:    conf.Client,
			Index:     conf.IndexName,
			BatchSize: 10,
			DocumentToFields: func(ctx context.Context, doc *schema.Document) (field2Value map[string]es8.FieldValue, err error) {
				var knowledgeName string
				if value, ok := ctx.Value(coretypes.KnowledgeName).(string); ok {
					knowledgeName = value
				} else {
					err = fmt.Errorf("必须提供知识库名称")
					return
				}
				if doc.MetaData != nil {
					// 存储ext数据
					marshal, _ := sonic.Marshal(getExtData(doc))
					doc.MetaData[coretypes.FieldExtra] = string(marshal)
				}
				return map[string]es8.FieldValue{
					coretypes.FieldContent: {
						Value:    doc.Content,
						EmbedKey: coretypes.FieldContentVector,
					},
					coretypes.FieldExtra: {
						Value: doc.MetaData[coretypes.FieldExtra],
					},
					coretypes.KnowledgeName: {
						Value: knowledgeName,
					},
					coretypes.FieldQAContent: {
						Value:    doc.MetaData[coretypes.FieldQAContent],
						EmbedKey: coretypes.FieldQAContentVector,
					},
				}, nil
			},
		}
		indexerConfig.Embedding = embeddingIns11
		idr, err = es8.NewIndexer(ctx, indexerConfig)
		if err != nil {
			return nil, err
		}
		return idr, nil
	} else {
		// Qdrant indexer
		idr, err = NewQdrantIndexer(ctx, &QdrantIndexerConfig{
			VectorStore: conf.VectorStore,
			IndexName:   conf.IndexName,
			Embedding:   embeddingIns11,
			BatchSize:   10,
			IsAsync:     true,
		})
		if err != nil {
			return nil, err
		}
		return idr, nil
	}
}
