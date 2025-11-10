package indexer

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
	coretypes "github.com/wangle201210/go-rag/server/core/types"
	"github.com/wangle201210/go-rag/server/core/vector"
)

// QdrantIndexerConfig Qdrant indexer 配置
type QdrantIndexerConfig struct {
	VectorStore vector.VectorStore
	IndexName   string
	Embedding   embedding.Embedder
	BatchSize   int
	IsAsync     bool // 是否异步模式（包含 QA 向量）
}

// QdrantIndexer Qdrant indexer 实现
type QdrantIndexer struct {
	config *QdrantIndexerConfig
}

// NewQdrantIndexer 创建 Qdrant indexer
func NewQdrantIndexer(ctx context.Context, config *QdrantIndexerConfig) (indexer.Indexer, error) {
	if config.BatchSize == 0 {
		config.BatchSize = 10
	}
	return &QdrantIndexer{
		config: config,
	}, nil
}

// Index 索引文档
func (idx *QdrantIndexer) Index(ctx context.Context, docs []*schema.Document, opts ...indexer.Option) ([]string, error) {
	if len(docs) == 0 {
		return nil, nil
	}

	var knowledgeName string
	if value, ok := ctx.Value(coretypes.KnowledgeName).(string); ok {
		knowledgeName = value
	} else {
		return nil, fmt.Errorf("必须提供知识库名称")
	}

	// 获取 Qdrant 客户端
	qdrantStore, ok := idx.config.VectorStore.(*vector.QdrantVectorStore)
	if !ok {
		return nil, fmt.Errorf("invalid vector store type")
	}

	var ids []string
	var points []*qdrant.PointStruct

	for _, doc := range docs {
		// 生成 ID
		if len(doc.ID) == 0 {
			doc.ID = uuid.New().String()
		}
		ids = append(ids, doc.ID)

		// 准备 payload
		payload := make(map[string]*qdrant.Value)
		payload[coretypes.FieldContent] = &qdrant.Value{
			Kind: &qdrant.Value_StringValue{StringValue: doc.Content},
		}
		payload[coretypes.KnowledgeName] = &qdrant.Value{
			Kind: &qdrant.Value_StringValue{StringValue: knowledgeName},
		}

		// 存储 ext 数据
		if doc.MetaData != nil {
			marshal, _ := sonic.Marshal(getExtData(doc))
			payload[coretypes.FieldExtra] = &qdrant.Value{
				Kind: &qdrant.Value_StringValue{StringValue: string(marshal)},
			}
		}

		// 生成 content 向量
		contentVec, err := idx.embedText(ctx, doc.Content)
		if err != nil {
			return nil, fmt.Errorf("failed to embed content: %w", err)
		}

		// 如果是异步模式，还需要生成 QA 向量
		if idx.config.IsAsync {
			if qaContent, ok := doc.MetaData[coretypes.FieldQAContent].(string); ok && qaContent != "" {
				payload[coretypes.FieldQAContent] = &qdrant.Value{
					Kind: &qdrant.Value_StringValue{StringValue: qaContent},
				}

				qaVec, err := idx.embedText(ctx, qaContent)
				if err != nil {
					return nil, fmt.Errorf("failed to embed qa content: %w", err)
				}

				// Qdrant 支持多向量，使用命名向量
				point := &qdrant.PointStruct{
					Id: &qdrant.PointId{
						PointIdOptions: &qdrant.PointId_Uuid{Uuid: doc.ID},
					},
					Vectors: &qdrant.Vectors{
						VectorsOptions: &qdrant.Vectors_Vectors{
							Vectors: &qdrant.NamedVectors{
								Vectors: map[string]*qdrant.Vector{
									coretypes.FieldContentVector:   {Data: contentVec},
									coretypes.FieldQAContentVector: {Data: qaVec},
								},
							},
						},
					},
					Payload: payload,
				}
				points = append(points, point)
			} else {
				// 没有 QA 内容，只使用 content 向量
				point := &qdrant.PointStruct{
					Id: &qdrant.PointId{
						PointIdOptions: &qdrant.PointId_Uuid{Uuid: doc.ID},
					},
					Vectors: &qdrant.Vectors{
						VectorsOptions: &qdrant.Vectors_Vector{
							Vector: &qdrant.Vector{Data: contentVec},
						},
					},
					Payload: payload,
				}
				points = append(points, point)
			}
		} else {
			// 非异步模式，只使用 content 向量
			point := &qdrant.PointStruct{
				Id: &qdrant.PointId{
					PointIdOptions: &qdrant.PointId_Uuid{Uuid: doc.ID},
				},
				Vectors: &qdrant.Vectors{
					VectorsOptions: &qdrant.Vectors_Vector{
						Vector: &qdrant.Vector{Data: contentVec},
					},
				},
				Payload: payload,
			}
			points = append(points, point)
		}
	}

	// 批量插入
	client := qdrantStore.GetClient()
	_, err := client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: idx.config.IndexName,
		Points:         points,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upsert points: %w", err)
	}

	return ids, nil
}

// embedText 生成文本向量
func (idx *QdrantIndexer) embedText(ctx context.Context, text string) ([]float32, error) {
	embedResp, err := idx.config.Embedding.EmbedStrings(ctx, []string{text})
	if err != nil {
		return nil, err
	}

	if len(embedResp) == 0 || len(embedResp[0]) == 0 {
		return nil, fmt.Errorf("empty embedding result")
	}

	// 转换为 float32
	vec := make([]float32, len(embedResp[0]))
	for i, v := range embedResp[0] {
		vec[i] = float32(v)
	}

	return vec, nil
}

// GetType 返回 indexer 类型
func (idx *QdrantIndexer) GetType() string {
	return "qdrant_indexer"
}

// Store 存储文档（实现 Indexer 接口）
func (idx *QdrantIndexer) Store(ctx context.Context, docs []*schema.Document, opts ...indexer.Option) ([]string, error) {
	return idx.Index(ctx, docs, opts...)
}
