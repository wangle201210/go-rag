package indexer

//
// type Orchestrator struct {
// 	docMapping service.IDocumentMapping
// 	docChunk   service.IDocumentChunk
// }
//
// func NewOrchestrator(docMapping service.IDocumentMapping, docChunk service.IDocumentChunk) *Orchestrator {
// 	return &Orchestrator{
// 		docMapping: docMapping,
// 		docChunk:   docChunk,
// 	}
// }
//
// func (o *Orchestrator) Orchestrate(ctx context.Context, doc *Document, chunks []*Chunk, knowledgeBaseID int64) error {
// 	// 创建文档映射关系
// 	docMapping := &model.DocumentMappingCreateReq{
// 		KnowledgeBaseId: knowledgeBaseID,
// 		DocumentId:      doc.ID,
// 		DocumentName:    doc.Name,
// 		DocumentType:    doc.Type,
// 		DocumentPath:    "", // 这里需要从外部传入
// 		DocumentSize:    int64(len(doc.Content)),
// 		ChunkCount:      len(chunks),
// 	}
//
// 	// 保存文档映射关系
// 	_, err := o.docMapping.Create(ctx, docMapping)
// 	if err != nil {
// 		return fmt.Errorf("failed to create document mapping: %v", err)
// 	}
//
// 	// 保存文档块
// 	for _, chunk := range chunks {
// 		// 将块元数据转换为JSON
// 		metadata, err := json.Marshal(chunk.Metadata)
// 		if err != nil {
// 			return fmt.Errorf("failed to marshal chunk metadata: %v", err)
// 		}
//
// 		// 创建文档块记录
// 		_, err = o.docChunk.Create(ctx, &model.DocumentChunkCreateReq{
// 			DocumentId: doc.ID,
// 			ChunkId:    chunk.ID,
// 			Content:    chunk.Content,
// 			Metadata:   string(metadata),
// 		})
// 		if err != nil {
// 			return fmt.Errorf("failed to create document chunk: %v", err)
// 		}
// 	}
//
// 	return nil
// }
//
// func (o *Orchestrator) Delete(ctx context.Context, documentID string) error {
// 	// 获取文档映射关系
// 	docMapping, err := o.docMapping.GetByDocumentId(ctx, documentID)
// 	if err != nil {
// 		return fmt.Errorf("failed to get document mapping: %v", err)
// 	}
//
// 	// 删除文档映射关系
// 	_, err = o.docMapping.Delete(ctx, &model.DocumentMappingDeleteReq{
// 		Id: docMapping.Id,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to delete document mapping: %v", err)
// 	}
//
// 	// 获取文档块
// 	chunks, err := o.docChunk.GetByDocumentId(ctx, documentID)
// 	if err != nil {
// 		return fmt.Errorf("failed to get document chunks: %v", err)
// 	}
//
// 	// 删除文档块
// 	for _, chunk := range chunks {
// 		_, err = o.docChunk.Delete(ctx, &model.DocumentChunkDeleteReq{
// 			Id: chunk.Id,
// 		})
// 		if err != nil {
// 			return fmt.Errorf("failed to delete document chunk: %v", err)
// 		}
// 	}
//
// 	return nil
// }
