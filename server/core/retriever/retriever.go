package retriever

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino-ext/components/retriever/es8"
	"github.com/cloudwego/eino-ext/components/retriever/es8/search_mode"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/wangle201210/go-rag/server/core/common"
	"github.com/wangle201210/go-rag/server/core/config"
	coretypes "github.com/wangle201210/go-rag/server/core/types"
)

// newRetriever component initialization function of node 'Retriever1' in graph 'retriever'
func newRetriever(ctx context.Context, conf *config.Config) (rtr retriever.Retriever, err error) {
	vectorField := coretypes.FieldContentVector
	if value, ok := ctx.Value(coretypes.RetrieverFieldKey).(string); ok {
		vectorField = value
	}

	embeddingIns, err := common.NewEmbedding(ctx, conf)
	if err != nil {
		return nil, err
	}

	// 根据客户端类型创建不同的 retriever
	if conf.Client != nil {
		// ES retriever
		retrieverConfig := &es8.RetrieverConfig{
			Client: conf.Client,
			Index:  conf.IndexName,
			SearchMode: search_mode.SearchModeDenseVectorSimilarity(
				search_mode.DenseVectorSimilarityTypeCosineSimilarity,
				vectorField,
			),
			ResultParser: EsHit2Document,
			Embedding:    embeddingIns,
		}
		rtr, err = es8.NewRetriever(ctx, retrieverConfig)
		if err != nil {
			return nil, err
		}
		return rtr, nil
	} else if conf.QdrantClient != nil {
		// Qdrant retriever
		rtr, err = NewQdrantRetriever(ctx, &QdrantRetrieverConfig{
			Client:      conf.QdrantClient,
			Collection:  conf.IndexName,
			Embedding:   embeddingIns,
			VectorField: vectorField,
		})
		if err != nil {
			return nil, err
		}
		return rtr, nil
	}

	return nil, fmt.Errorf("no valid client configuration found")
}

func EsHit2Document(ctx context.Context, hit types.Hit) (doc *schema.Document, err error) {
	doc = &schema.Document{
		ID:       *hit.Id_,
		MetaData: map[string]any{},
	}

	var src map[string]any
	if err = sonic.Unmarshal(hit.Source_, &src); err != nil {
		return nil, err
	}

	for field, val := range src {
		switch field {
		case coretypes.FieldContent:
			doc.Content = val.(string)
		case coretypes.FieldContentVector:
			var v []float64
			for _, item := range val.([]interface{}) {
				v = append(v, item.(float64))
			}
			doc.WithDenseVector(v)
		case coretypes.FieldQAContentVector, coretypes.FieldQAContent:
			// 这两个字段都不返回

		case coretypes.FieldExtra:
			if val == nil {
				continue
			}
			doc.MetaData[coretypes.FieldExtra] = val.(string)
		case coretypes.KnowledgeName:
			doc.MetaData[coretypes.KnowledgeName] = val.(string)
		default:
			return nil, fmt.Errorf("unexpected field=%s, val=%v", field, val)
		}
	}

	if hit.Score_ != nil {
		doc.WithScore(float64(*hit.Score_))
	}

	return doc, nil
}
