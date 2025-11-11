package indexer

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
	coretypes "github.com/wangle201210/go-rag/server/core/types"
)

// StoreWithNamedVectors 使用命名向量存储文档到 Qdrant
// 这是自定义实现，因为 eino-ext 的 indexer 不支持命名向量
func (idx *QdrantIndexer) StoreWithNamedVectors(ctx context.Context, docs []*schema.Document, opts ...indexer.Option) ([]string, error) {
	if len(docs) == 0 {
		return nil, nil
	}

	var knowledgeName string
	if value, ok := ctx.Value(coretypes.KnowledgeName).(string); ok {
		knowledgeName = value
	} else {
		return nil, fmt.Errorf("必须提供知识库名称")
	}

	g.Log().Infof(ctx, "QdrantIndexer.StoreWithNamedVectors: storing %d documents to collection %s, knowledge_name=%s", len(docs), idx.config.Collection, knowledgeName)

	// 准备 points
	points := make([]*qdrant.PointStruct, 0, len(docs))
	ids := make([]string, 0, len(docs))

	for _, doc := range docs {
		// 生成 ID
		if len(doc.ID) == 0 {
			doc.ID = uuid.New().String()
		}
		ids = append(ids, doc.ID)

		// 生成 embedding
		embeddings, err := idx.config.Embedding.EmbedStrings(ctx, []string{doc.Content})
		if err != nil {
			g.Log().Errorf(ctx, "Failed to embed document %s: %v", doc.ID, err)
			return nil, fmt.Errorf("failed to embed document: %w", err)
		}
		if len(embeddings) == 0 {
			return nil, fmt.Errorf("embedding returned empty result")
		}

		// 转换为 float32
		vec32 := make([]float32, len(embeddings[0]))
		for i, v := range embeddings[0] {
			vec32[i] = float32(v)
		}

		// 准备 payload
		payload := make(map[string]*qdrant.Value)
		payload[coretypes.FieldContent] = &qdrant.Value{
			Kind: &qdrant.Value_StringValue{StringValue: doc.Content},
		}
		payload[coretypes.KnowledgeName] = &qdrant.Value{
			Kind: &qdrant.Value_StringValue{StringValue: knowledgeName},
		}

		// 添加额外的 metadata
		if doc.MetaData != nil {
			extData := getExtData(doc)
			if len(extData) > 0 {
				marshal, _ := sonic.Marshal(extData)
				payload[coretypes.FieldExtra] = &qdrant.Value{
					Kind: &qdrant.Value_StringValue{StringValue: string(marshal)},
				}
			}
		}

		// 创建命名向量（只存储 content_vector，qa_content_vector 由异步任务处理）
		vectors := &qdrant.Vectors{
			VectorsOptions: &qdrant.Vectors_Vectors{
				Vectors: &qdrant.NamedVectors{
					Vectors: map[string]*qdrant.Vector{
						coretypes.FieldContentVector: {
							Data: vec32,
						},
					},
				},
			},
		}

		// 创建 point
		point := &qdrant.PointStruct{
			Id: &qdrant.PointId{
				PointIdOptions: &qdrant.PointId_Uuid{Uuid: doc.ID},
			},
			Vectors: vectors,
			Payload: payload,
		}

		points = append(points, point)
	}

	// 批量存储到 Qdrant
	_, err := idx.config.Client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: idx.config.Collection,
		Points:         points,
	})
	if err != nil {
		g.Log().Errorf(ctx, "QdrantIndexer.StoreWithNamedVectors failed: %v", err)
		return nil, fmt.Errorf("failed to upsert points: %w", err)
	}

	g.Log().Infof(ctx, "QdrantIndexer.StoreWithNamedVectors success: stored %d documents, IDs: %v", len(ids), ids)

	return ids, nil
}
