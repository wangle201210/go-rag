package server

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/wangle201210/go-rag/server/core/model"
)

type IDocumentMapping interface {
	Create(ctx context.Context, req *model.DocumentMappingCreateReq) (*model.DocumentMappingCreateRes, error)
	Update(ctx context.Context, req *model.DocumentMappingUpdateReq) (*model.DocumentMappingUpdateRes, error)
	Delete(ctx context.Context, req *model.DocumentMappingDeleteReq) (*model.DocumentMappingDeleteRes, error)
	GetList(ctx context.Context, req *model.DocumentMappingListReq) (*model.DocumentMappingListRes, error)
	GetByDocumentId(ctx context.Context, documentId string) (*model.DocumentMapping, error)
}

type documentMappingImpl struct{}

func NewDocumentMapping() IDocumentMapping {
	return &documentMappingImpl{}
}

func (s *documentMappingImpl) Create(ctx context.Context, req *model.DocumentMappingCreateReq) (*model.DocumentMappingCreateRes, error) {
	result, err := g.DB().Model("document_mapping").Insert(g.Map{
		"knowledge_base_id": req.KnowledgeBaseId,
		"document_id":       req.DocumentId,
		"document_name":     req.DocumentName,
		"document_type":     req.DocumentType,
		"document_path":     req.DocumentPath,
		"document_size":     req.DocumentSize,
		"chunk_count":       req.ChunkCount,
		"status":            1,
	})
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &model.DocumentMappingCreateRes{
		Id: id,
	}, nil
}

func (s *documentMappingImpl) Update(ctx context.Context, req *model.DocumentMappingUpdateReq) (*model.DocumentMappingUpdateRes, error) {
	_, err := g.DB().Model("document_mapping").
		Where("id", req.Id).
		Update(g.Map{
			"knowledge_base_id": req.KnowledgeBaseId,
			"document_name":     req.DocumentName,
			"document_type":     req.DocumentType,
			"document_path":     req.DocumentPath,
			"document_size":     req.DocumentSize,
			"chunk_count":       req.ChunkCount,
			"status":            req.Status,
		})
	if err != nil {
		return nil, err
	}

	return &model.DocumentMappingUpdateRes{
		Success: true,
	}, nil
}

func (s *documentMappingImpl) Delete(ctx context.Context, req *model.DocumentMappingDeleteReq) (*model.DocumentMappingDeleteRes, error) {
	_, err := g.DB().Model("document_mapping").Where("id", req.Id).Delete()
	if err != nil {
		return nil, err
	}

	return &model.DocumentMappingDeleteRes{
		Success: true,
	}, nil
}

func (s *documentMappingImpl) GetList(ctx context.Context, req *model.DocumentMappingListReq) (*model.DocumentMappingListRes, error) {
	m := g.DB().Model("document_mapping")

	if req.KnowledgeBaseId > 0 {
		m = m.Where("knowledge_base_id", req.KnowledgeBaseId)
	}
	if req.Status > 0 {
		m = m.Where("status", req.Status)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	list, err := m.Page(req.Page, req.PageSize).OrderDesc("id").All()
	if err != nil {
		return nil, err
	}

	var items []model.DocumentMapping
	if err := list.Structs(&items); err != nil {
		return nil, err
	}

	return &model.DocumentMappingListRes{
		List:  items,
		Total: total,
	}, nil
}

func (s *documentMappingImpl) GetByDocumentId(ctx context.Context, documentId string) (*model.DocumentMapping, error) {
	var mapping model.DocumentMapping
	err := g.DB().Model("document_mapping").
		Where("document_id", documentId).
		Scan(&mapping)
	if err != nil {
		return nil, err
	}
	return &mapping, nil
}
