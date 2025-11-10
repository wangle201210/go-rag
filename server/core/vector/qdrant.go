package vector

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/qdrant/go-client/qdrant"
	coretypes "github.com/wangle201210/go-rag/server/core/types"
)

// QdrantVectorStore Qdrant 向量存储实现
type QdrantVectorStore struct {
	client *qdrant.Client
}

// NewQdrantVectorStore 创建 Qdrant 向量存储实例
func NewQdrantVectorStore(cfg *QdrantConfig) (*QdrantVectorStore, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   cfg.Address,
		Port:   cfg.Port,
		APIKey: cfg.APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create qdrant client: %w", err)
	}

	return &QdrantVectorStore{
		client: client,
	}, nil
}

// GetClient 获取 Qdrant 客户端（用于兼容现有代码）
func (s *QdrantVectorStore) GetClient() *qdrant.Client {
	return s.client
}

// CreateIndex 创建集合
func (s *QdrantVectorStore) CreateIndex(ctx context.Context, indexName string) error {
	// 创建集合，包含两个向量字段：content_vector 和 qa_content_vector
	err := s.client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: indexName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     1024,
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	// 创建 payload 索引以支持过滤
	_, err = s.client.CreateFieldIndex(ctx, &qdrant.CreateFieldIndexCollection{
		CollectionName: indexName,
		FieldName:      coretypes.KnowledgeName,
		FieldType:      qdrant.FieldType_FieldTypeKeyword.Enum(),
	})
	if err != nil {
		return fmt.Errorf("failed to create field index: %w", err)
	}

	return nil
}

// IndexExists 检查集合是否存在
func (s *QdrantVectorStore) IndexExists(ctx context.Context, indexName string) (bool, error) {
	exists, err := s.client.CollectionExists(ctx, indexName)
	if err != nil {
		return false, fmt.Errorf("failed to check collection exists: %w", err)
	}
	return exists, nil
}

// DeleteDocument 删除文档
func (s *QdrantVectorStore) DeleteDocument(ctx context.Context, indexName, documentID string) error {
	_, err := s.client.Delete(ctx, &qdrant.DeletePoints{
		CollectionName: indexName,
		Points: &qdrant.PointsSelector{
			PointsSelectorOneOf: &qdrant.PointsSelector_Points{
				Points: &qdrant.PointsIdsList{
					Ids: []*qdrant.PointId{
						{PointIdOptions: &qdrant.PointId_Uuid{Uuid: documentID}},
					},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	return nil
}

// GetKnowledgeBaseList 获取知识库列表
func (s *QdrantVectorStore) GetKnowledgeBaseList(ctx context.Context, indexName string) ([]string, error) {
	// Qdrant 不直接支持聚合，需要通过滚动查询获取所有唯一的知识库名称
	// 这里使用一个简化的实现：查询所有点并去重
	knowledgeMap := make(map[string]bool)

	offset := (*qdrant.PointId)(nil)
	limit := uint32(100)

	for {
		scrollResp, err := s.client.Scroll(ctx, &qdrant.ScrollPoints{
			CollectionName: indexName,
			Limit:          &limit,
			Offset:         offset,
			WithPayload:    &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to scroll points: %w", err)
		}

		if len(scrollResp) == 0 {
			break
		}

		for _, point := range scrollResp {
			if payload := point.GetPayload(); payload != nil {
				if knowledgeName, ok := payload[coretypes.KnowledgeName]; ok {
					if name, ok := knowledgeName.GetKind().(*qdrant.Value_StringValue); ok {
						knowledgeMap[name.StringValue] = true
					}
				}
			}
		}

		if len(scrollResp) < int(limit) {
			break
		}
		offset = scrollResp[len(scrollResp)-1].Id
	}

	var list []string
	for name := range knowledgeMap {
		list = append(list, name)
	}

	return list, nil
}

// SearchDocuments 搜索文档
func (s *QdrantVectorStore) SearchDocuments(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	// 构建过滤条件
	var filter *qdrant.Filter
	if req.KnowledgeName != "" {
		filter = &qdrant.Filter{
			Must: []*qdrant.Condition{
				{
					ConditionOneOf: &qdrant.Condition_Field{
						Field: &qdrant.FieldCondition{
							Key: coretypes.KnowledgeName,
							Match: &qdrant.Match{
								MatchValue: &qdrant.Match_Keyword{
									Keyword: req.KnowledgeName,
								},
							},
						},
					},
				},
			},
		}
	}

	if len(req.DocIDs) > 0 {
		ids := make([]*qdrant.PointId, len(req.DocIDs))
		for i, id := range req.DocIDs {
			ids[i] = &qdrant.PointId{PointIdOptions: &qdrant.PointId_Uuid{Uuid: id}}
		}

		hasFilter := &qdrant.Condition{
			ConditionOneOf: &qdrant.Condition_HasId{
				HasId: &qdrant.HasIdCondition{
					HasId: ids,
				},
			},
		}

		if filter == nil {
			filter = &qdrant.Filter{Must: []*qdrant.Condition{hasFilter}}
		} else {
			filter.Must = append(filter.Must, hasFilter)
		}
	}

	// 如果有向量查询，使用 Query API
	if qdrantReq, ok := req.Query.(*QdrantSearchQuery); ok && qdrantReq != nil {
		searchResp, err := s.client.Query(ctx, &qdrant.QueryPoints{
			CollectionName: qdrantReq.CollectionName,
			Query:          qdrantReq.Query,
			Filter:         filter,
			Limit:          uint64ptr(uint64(req.Size)),
			WithPayload:    &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
			WithVectors:    &qdrant.WithVectorsSelector{SelectorOptions: &qdrant.WithVectorsSelector_Enable{Enable: true}},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to search: %w", err)
		}

		var docs []*schema.Document
		for _, point := range searchResp {
			doc, err := s.qdrantPoint2Document(ctx, point)
			if err != nil {
				g.Log().Errorf(ctx, "qdrantPoint2Document failed, err=%v", err)
				continue
			}
			docs = append(docs, doc)
		}

		return &SearchResponse{
			Documents: docs,
			Total:     int64(len(docs)),
		}, nil
	}

	// 否则使用 Scroll API（用于按 ID 和过滤条件获取文档）
	scrollResp, err := s.client.Scroll(ctx, &qdrant.ScrollPoints{
		CollectionName: req.IndexName,
		Filter:         filter,
		Limit:          uint32ptr(uint32(req.Size)),
		WithPayload:    &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
		WithVectors:    &qdrant.WithVectorsSelector{SelectorOptions: &qdrant.WithVectorsSelector_Enable{Enable: true}},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scroll: %w", err)
	}

	var docs []*schema.Document
	for _, point := range scrollResp {
		doc, err := s.qdrantScrollPoint2Document(ctx, point)
		if err != nil {
			g.Log().Errorf(ctx, "qdrantScrollPoint2Document failed, err=%v", err)
			continue
		}
		docs = append(docs, doc)
	}

	return &SearchResponse{
		Documents: docs,
		Total:     int64(len(docs)),
	}, nil
}

// qdrantPoint2Document 将 Qdrant ScoredPoint 转换为 Document
func (s *QdrantVectorStore) qdrantPoint2Document(_ context.Context, point *qdrant.ScoredPoint) (*schema.Document, error) {
	var docID string
	switch id := point.Id.GetPointIdOptions().(type) {
	case *qdrant.PointId_Uuid:
		docID = id.Uuid
	case *qdrant.PointId_Num:
		docID = fmt.Sprintf("%d", id.Num)
	default:
		return nil, fmt.Errorf("unsupported point id type")
	}

	doc := &schema.Document{
		ID:       docID,
		MetaData: map[string]any{},
	}

	payload := point.GetPayload()
	if payload != nil {
		// 提取 content
		if content, ok := payload[coretypes.FieldContent]; ok {
			if val, ok := content.GetKind().(*qdrant.Value_StringValue); ok {
				doc.Content = val.StringValue
			}
		}

		// 提取 extra
		if extra, ok := payload[coretypes.FieldExtra]; ok {
			if val, ok := extra.GetKind().(*qdrant.Value_StringValue); ok {
				doc.MetaData[coretypes.FieldExtra] = val.StringValue
			}
		}

		// 提取 knowledge_name
		if knowledgeName, ok := payload[coretypes.KnowledgeName]; ok {
			if val, ok := knowledgeName.GetKind().(*qdrant.Value_StringValue); ok {
				doc.MetaData[coretypes.KnowledgeName] = val.StringValue
			}
		}
	}

	// 提取向量
	if vectors := point.GetVectors(); vectors != nil {
		switch v := vectors.GetVectorsOptions().(type) {
		case *qdrant.VectorsOutput_Vector:
			// 转换 float32 到 float64
			data := v.Vector.GetData()
			vec64 := make([]float64, len(data))
			for i, val := range data {
				vec64[i] = float64(val)
			}
			doc.WithDenseVector(vec64)
		}
	}

	// 提取分数
	doc.WithScore(float64(point.Score))

	return doc, nil
}

// qdrantScrollPoint2Document 将 Qdrant RetrievedPoint 转换为 Document
func (s *QdrantVectorStore) qdrantScrollPoint2Document(_ context.Context, point *qdrant.RetrievedPoint) (*schema.Document, error) {
	var docID string
	switch id := point.Id.GetPointIdOptions().(type) {
	case *qdrant.PointId_Uuid:
		docID = id.Uuid
	case *qdrant.PointId_Num:
		docID = fmt.Sprintf("%d", id.Num)
	default:
		return nil, fmt.Errorf("unsupported point id type")
	}

	doc := &schema.Document{
		ID:       docID,
		MetaData: map[string]any{},
	}

	payload := point.GetPayload()
	if payload != nil {
		// 提取 content
		if content, ok := payload[coretypes.FieldContent]; ok {
			if val, ok := content.GetKind().(*qdrant.Value_StringValue); ok {
				doc.Content = val.StringValue
			}
		}

		// 提取 extra
		if extra, ok := payload[coretypes.FieldExtra]; ok {
			if val, ok := extra.GetKind().(*qdrant.Value_StringValue); ok {
				doc.MetaData[coretypes.FieldExtra] = val.StringValue
			}
		}

		// 提取 knowledge_name
		if knowledgeName, ok := payload[coretypes.KnowledgeName]; ok {
			if val, ok := knowledgeName.GetKind().(*qdrant.Value_StringValue); ok {
				doc.MetaData[coretypes.KnowledgeName] = val.StringValue
			}
		}
	}

	// 提取向量
	if vectors := point.GetVectors(); vectors != nil {
		switch v := vectors.GetVectorsOptions().(type) {
		case *qdrant.VectorsOutput_Vector:
			// 转换 float32 到 float64
			data := v.Vector.GetData()
			vec64 := make([]float64, len(data))
			for i, val := range data {
				vec64[i] = float64(val)
			}
			doc.WithDenseVector(vec64)
		}
	}

	return doc, nil
}

// Close 关闭连接
func (s *QdrantVectorStore) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

// uint32ptr 辅助函数
func uint32ptr(v uint32) *uint32 {
	return &v
}

// uint64ptr 辅助函数
func uint64ptr(v uint64) *uint64 {
	return &v
}

// QdrantSearchQuery Qdrant 搜索查询
type QdrantSearchQuery struct {
	CollectionName string
	Query          *qdrant.Query
}
