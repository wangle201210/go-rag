package indexer

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/components/indexer/es8"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/wangle201210/go-rag/server/core/common"
	"github.com/wangle201210/go-rag/server/core/config"
	coretypes "github.com/wangle201210/go-rag/server/core/types"
)

// newIndexer component initialization function of node 'Indexer2' in graph 'rag'
func newIndexer(ctx context.Context, conf *config.Config) (idr indexer.Indexer, err error) {
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
				// 没传入才需要生成
				if len(doc.ID) == 0 {
					doc.ID = uuid.New().String()
				}
				if doc.MetaData != nil {
					// 存储ext数据
					marshal, _ := sonic.Marshal(getExtData(doc))
					doc.MetaData[coretypes.FieldExtra] = string(marshal)
				}
				return map[string]es8.FieldValue{
					coretypes.FieldContent: {
						Value:    doc.Content,
						EmbedKey: coretypes.FieldContentVector, // 这里也可以考虑注释掉，等后续异步执行
					},
					coretypes.FieldExtra: {
						Value: doc.MetaData[coretypes.FieldExtra],
					},
					coretypes.KnowledgeName: {
						Value: knowledgeName,
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
			IsAsync:     false,
		})
		if err != nil {
			return nil, err
		}
		return idr, nil
	}
}

func getExtData(doc *schema.Document) map[string]any {
	if doc.MetaData == nil {
		return nil
	}
	res := make(map[string]any)
	for _, key := range coretypes.ExtKeys {
		if v, e := doc.MetaData[key]; e {
			res[key] = v
		}
	}
	return res
}
