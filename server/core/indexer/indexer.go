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
)

// newIndexer component initialization function of node 'Indexer2' in graph 'rag'
func newIndexer(ctx context.Context, conf *config.Config) (idr indexer.Indexer, err error) {
	indexerConfig := &es8.IndexerConfig{
		Client:    conf.Client,
		Index:     conf.IndexName,
		BatchSize: 10,
		DocumentToFields: func(ctx context.Context, doc *schema.Document) (field2Value map[string]es8.FieldValue, err error) {
			var knowledgeName string
			if value, ok := ctx.Value(common.KnowledgeName).(string); ok {
				knowledgeName = value
			} else {
				err = fmt.Errorf("必须提供知识库名称")
				return
			}
			doc.ID = uuid.New().String()
			if doc.MetaData != nil {
				marshal, _ := sonic.Marshal(doc.MetaData)
				doc.MetaData[common.DocExtra] = string(marshal)
			}
			return map[string]es8.FieldValue{
				common.FieldContent: {
					Value:    getMdContentWithTitle(doc),
					EmbedKey: common.FieldContentVector, // vectorize doc content and save vector to field "content_vector"
				},
				common.FieldExtra: {
					Value: doc.MetaData[common.DocExtra],
				},
				common.KnowledgeName: {
					Value: knowledgeName,
				},
			}, nil
		},
	}
	embeddingIns11, err := common.NewEmbedding(ctx, conf)
	if err != nil {
		return nil, err
	}
	indexerConfig.Embedding = embeddingIns11
	idr, err = es8.NewIndexer(ctx, indexerConfig)
	if err != nil {
		return nil, err
	}
	return idr, nil
}

func getMdContentWithTitle(doc *schema.Document) string {
	if doc.MetaData == nil {
		return doc.Content
	}
	title := ""
	list := []string{"h1", "h2", "h3", "h4", "h5", "h6"}
	for _, v := range list {
		if d, e := doc.MetaData[v].(string); e && len(d) > 0 {
			title += fmt.Sprintf("%s:%s ", v, d)
		}
	}
	if len(title) == 0 {
		return doc.Content
	}
	return title + "\n" + doc.Content
}

//
// type Indexer struct {
// 	loader       *Loader
// 	parser       *Parser
// 	transformer  *Transformer
// 	orchestrator *Orchestrator
// 	docMapping   service.IDocumentMapping
// }
//
// func NewIndexer(docMapping service.IDocumentMapping) *Indexer {
// 	loader := NewLoader()
// 	parser := NewParser()
// 	transformer := NewTransformer()
// 	orchestrator := NewOrchestrator()
//
// 	return &Indexer{
// 		loader:       loader,
// 		parser:       parser,
// 		transformer:  transformer,
// 		orchestrator: orchestrator,
// 		docMapping:   docMapping,
// 	}
// }
//
// func (i *Indexer) Index(ctx context.Context, knowledgeBaseId int64, filePath string) error {
// 	// 1. 加载文档
// 	content, err := i.loader.Load(filePath)
// 	if err != nil {
// 		return fmt.Errorf("加载文档失败: %v", err)
// 	}
//
// 	// 2. 解析文档
// 	doc, err := i.parser.Parse(content)
// 	if err != nil {
// 		return fmt.Errorf("解析文档失败: %v", err)
// 	}
//
// 	// 3. 转换文档
// 	chunks, err := i.transformer.Transform(doc)
// 	if err != nil {
// 		return fmt.Errorf("转换文档失败: %v", err)
// 	}
//
// 	// 4. 创建文档映射
// 	mapping := &model.DocumentMappingCreateReq{
// 		KnowledgeBaseId: knowledgeBaseId,
// 		DocumentId:      doc.ID,
// 		DocumentName:    doc.Name,
// 		DocumentType:    doc.Type,
// 		DocumentPath:    filePath,
// 		DocumentSize:    int64(len(content)),
// 		ChunkCount:      len(chunks),
// 	}
//
// 	_, err = i.docMapping.Create(ctx, mapping)
// 	if err != nil {
// 		return fmt.Errorf("创建文档映射失败: %v", err)
// 	}
//
// 	// 5. 编排和存储
// 	err = i.orchestrator.Process(ctx, chunks)
// 	if err != nil {
// 		return fmt.Errorf("处理文档块失败: %v", err)
// 	}
//
// 	return nil
// }
//
// func (i *Indexer) Delete(ctx context.Context, documentId string) error {
// 	// 1. 获取文档映射
// 	mapping, err := i.docMapping.GetByDocumentId(ctx, documentId)
// 	if err != nil {
// 		return fmt.Errorf("获取文档映射失败: %v", err)
// 	}
//
// 	// 2. 删除文档映射
// 	_, err = i.docMapping.Delete(ctx, &model.DocumentMappingDeleteReq{
// 		Id: mapping.Id,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("删除文档映射失败: %v", err)
// 	}
//
// 	// 3. 删除文档块
// 	err = i.orchestrator.Delete(ctx, documentId)
// 	if err != nil {
// 		return fmt.Errorf("删除文档块失败: %v", err)
// 	}
//
// 	return nil
// }
