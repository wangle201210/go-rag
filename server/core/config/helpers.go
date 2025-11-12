package config

import (
	"context"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/schema"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/exists"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/qdrant/go-client/qdrant"
	coretypes "github.com/wangle201210/go-rag/server/core/types"
)

// IndexExists 检查索引/集合是否存在
func (c *Config) IndexExists(ctx context.Context) (bool, error) {
	if c.Client != nil {
		// ES
		return exists.NewExistsFunc(c.Client)(c.IndexName).Do(ctx)
	} else if c.QdrantClient != nil {
		// Qdrant
		return c.QdrantClient.CollectionExists(ctx, c.IndexName)
	}
	return false, fmt.Errorf("no valid client configuration")
}

// CreateIndex 创建索引/集合
func (c *Config) CreateIndex(ctx context.Context) error {
	if c.Client != nil {
		// ES
		_, err := create.NewCreateFunc(c.Client)(c.IndexName).Request(&create.Request{
			Mappings: &types.TypeMapping{
				Properties: map[string]types.Property{
					coretypes.FieldContent:  types.NewTextProperty(),
					coretypes.FieldExtra:    types.NewTextProperty(),
					coretypes.KnowledgeName: types.NewKeywordProperty(),
					coretypes.FieldContentVector: &types.DenseVectorProperty{
						Dims:       Of(1024),
						Index:      Of(true),
						Similarity: Of("cosine"),
					},
					coretypes.FieldQAContentVector: &types.DenseVectorProperty{
						Dims:       Of(1024),
						Index:      Of(true),
						Similarity: Of("cosine"),
					},
				},
			},
		}).Do(ctx)
		return err
	} else if c.QdrantClient != nil {
		// Qdrant - 创建集合，支持命名向量
		vectorsMap := map[string]*qdrant.VectorParams{
			coretypes.FieldContentVector: {
				Size:     1024,
				Distance: qdrant.Distance_Cosine,
			},
			coretypes.FieldQAContentVector: {
				Size:     1024,
				Distance: qdrant.Distance_Cosine,
			},
		}
		err := c.QdrantClient.CreateCollection(ctx, &qdrant.CreateCollection{
			CollectionName: c.IndexName,
			VectorsConfig:  qdrant.NewVectorsConfigMap(vectorsMap),
		})
		if err != nil {
			return fmt.Errorf("failed to create collection: %w", err)
		}

		// 创建 payload 索引以支持过滤
		_, err = c.QdrantClient.CreateFieldIndex(ctx, &qdrant.CreateFieldIndexCollection{
			CollectionName: c.IndexName,
			FieldName:      coretypes.KnowledgeName,
			FieldType:      qdrant.FieldType_FieldTypeKeyword.Enum(),
		})
		if err != nil {
			return fmt.Errorf("failed to create field index: %w", err)
		}

		return nil
	}
	return fmt.Errorf("no valid client configuration")
}

// DeleteDocument 删除文档
func (c *Config) DeleteDocument(ctx context.Context, documentID string) error {
	if c.Client != nil {
		// ES
		res, err := c.Client.Delete(c.IndexName, documentID)
		if err != nil {
			return fmt.Errorf("delete document failed: %w", err)
		}
		defer res.Body.Close()

		if res.IsError() {
			return fmt.Errorf("delete document failed: %s", res.String())
		}
		return nil
	} else if c.QdrantClient != nil {
		// Qdrant
		_, err := c.QdrantClient.Delete(ctx, &qdrant.DeletePoints{
			CollectionName: c.IndexName,
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
	return fmt.Errorf("no valid client configuration")
}

// GetKnowledgeBaseList 获取知识库列表
func (c *Config) GetKnowledgeBaseList(ctx context.Context) ([]string, error) {
	if c.Client != nil {
		// ES - 使用聚合查询
		// TODO: 实现 ES 的知识库列表获取
		return nil, fmt.Errorf("ES GetKnowledgeBaseList not implemented yet")
	} else if c.QdrantClient != nil {
		// Qdrant - 滚动查询并去重
		knowledgeMap := make(map[string]bool)

		offset := (*qdrant.PointId)(nil)
		limit := uint32(100)

		for {
			scrollResp, err := c.QdrantClient.Scroll(ctx, &qdrant.ScrollPoints{
				CollectionName: c.IndexName,
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
	return nil, fmt.Errorf("no valid client configuration")
}

// SearchDocumentsByIDs 根据文档 ID 列表搜索文档
func (c *Config) SearchDocumentsByIDs(ctx context.Context, knowledgeName string, docIDs []string, size int) ([]*schema.Document, error) {
	if c.Client != nil {
		// ES
		esQuery := &types.Query{
			Bool: &types.BoolQuery{
				Must: []types.Query{
					{Match: map[string]types.MatchQuery{coretypes.KnowledgeName: {Query: knowledgeName}}},
					{Terms: &types.TermsQuery{TermsQuery: map[string]types.TermsQueryField{"_id": docIDs}}},
				},
			},
		}

		sreq := search.NewRequest()
		sreq.Query = esQuery
		sreq.Size = Of(size)

		searchResp, err := search.NewSearchFunc(c.Client)().
			Index(c.IndexName).
			Request(sreq).
			Do(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to search: %w", err)
		}

		var docs []*schema.Document
		for _, hit := range searchResp.Hits.Hits {
			doc, err := esHit2Document(ctx, hit)
			if err != nil {
				continue
			}
			docs = append(docs, doc)
		}

		return docs, nil
	} else if c.QdrantClient != nil {
		// Qdrant
		ids := make([]*qdrant.PointId, len(docIDs))
		for i, id := range docIDs {
			ids[i] = &qdrant.PointId{PointIdOptions: &qdrant.PointId_Uuid{Uuid: id}}
		}

		filter := &qdrant.Filter{
			Must: []*qdrant.Condition{
				{
					ConditionOneOf: &qdrant.Condition_Field{
						Field: &qdrant.FieldCondition{
							Key: coretypes.KnowledgeName,
							Match: &qdrant.Match{
								MatchValue: &qdrant.Match_Keyword{
									Keyword: knowledgeName,
								},
							},
						},
					},
				},
				{
					ConditionOneOf: &qdrant.Condition_HasId{
						HasId: &qdrant.HasIdCondition{
							HasId: ids,
						},
					},
				},
			},
		}

		limit := uint32(size)
		scrollResp, err := c.QdrantClient.Scroll(ctx, &qdrant.ScrollPoints{
			CollectionName: c.IndexName,
			Filter:         filter,
			Limit:          &limit,
			WithPayload:    &qdrant.WithPayloadSelector{SelectorOptions: &qdrant.WithPayloadSelector_Enable{Enable: true}},
			WithVectors:    &qdrant.WithVectorsSelector{SelectorOptions: &qdrant.WithVectorsSelector_Enable{Enable: true}},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to scroll: %w", err)
		}

		var docs []*schema.Document
		for _, point := range scrollResp {
			doc, err := qdrantPoint2Document(ctx, point)
			if err != nil {
				continue
			}
			docs = append(docs, doc)
		}

		return docs, nil
	}
	return nil, fmt.Errorf("no valid client configuration")
}

// esHit2Document 将 ES Hit 转换为 Document
func esHit2Document(ctx context.Context, hit types.Hit) (*schema.Document, error) {
	doc := &schema.Document{
		ID:       *hit.Id_,
		MetaData: map[string]any{},
	}

	var src map[string]any
	if err := sonic.Unmarshal(hit.Source_, &src); err != nil {
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

// qdrantPoint2Document 将 Qdrant RetrievedPoint 转换为 Document
func qdrantPoint2Document(_ context.Context, point *qdrant.RetrievedPoint) (*schema.Document, error) {
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

// Of 辅助函数
func Of[T any](v T) *T {
	return &v
}
