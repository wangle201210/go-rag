package server

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/wangle201210/go-rag/server/core/model"
)

type IDocumentChunk interface {
	Create(ctx context.Context, req *model.DocumentChunkCreateReq) (*model.DocumentChunkCreateRes, error)
	Update(ctx context.Context, req *model.DocumentChunkUpdateReq) (*model.DocumentChunkUpdateRes, error)
	Delete(ctx context.Context, req *model.DocumentChunkDeleteReq) (*model.DocumentChunkDeleteRes, error)
	GetList(ctx context.Context, req *model.DocumentChunkListReq) (*model.DocumentChunkListRes, error)
	GetByDocumentId(ctx context.Context, documentId string) ([]*model.DocumentChunk, error)
}

type documentChunkImpl struct{}

func NewDocumentChunk() IDocumentChunk {
	return &documentChunkImpl{}
}

func (s *documentChunkImpl) Create(ctx context.Context, req *model.DocumentChunkCreateReq) (*model.DocumentChunkCreateRes, error) {
	result, err := g.DB().Model("document_chunk").Insert(g.Map{
		"document_id": req.DocumentId,
		"chunk_id":    req.ChunkId,
		"content":     req.Content,
		"metadata":    req.Metadata,
	})
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &model.DocumentChunkCreateRes{
		Id: id,
	}, nil
}

func (s *documentChunkImpl) Update(ctx context.Context, req *model.DocumentChunkUpdateReq) (*model.DocumentChunkUpdateRes, error) {
	updates := g.Map{}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Metadata != "" {
		updates["metadata"] = req.Metadata
	}

	_, err := g.DB().Model("document_chunk").
		Where("id", req.Id).
		Update(updates)
	if err != nil {
		return nil, err
	}

	return &model.DocumentChunkUpdateRes{
		Id: req.Id,
	}, nil
}

func (s *documentChunkImpl) Delete(ctx context.Context, req *model.DocumentChunkDeleteReq) (*model.DocumentChunkDeleteRes, error) {
	_, err := g.DB().Model("document_chunk").Where("id", req.Id).Delete()
	if err != nil {
		return nil, err
	}

	return &model.DocumentChunkDeleteRes{
		Id: req.Id,
	}, nil
}

func (s *documentChunkImpl) GetList(ctx context.Context, req *model.DocumentChunkListReq) (*model.DocumentChunkListRes, error) {
	m := g.DB().Model("document_chunk").Where("document_id", req.DocumentId)

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	list, err := m.Page(req.Page, req.PageSize).OrderDesc("id").All()
	if err != nil {
		return nil, err
	}

	var items []model.DocumentChunk
	if err := list.Structs(&items); err != nil {
		return nil, err
	}

	return &model.DocumentChunkListRes{
		Total: int64(total),
		List:  items,
	}, nil
}

func (s *documentChunkImpl) GetByDocumentId(ctx context.Context, documentId string) ([]*model.DocumentChunk, error) {
	list, err := g.DB().Model("document_chunk").
		Where("document_id", documentId).
		OrderAsc("id").
		All()
	if err != nil {
		return nil, err
	}

	var items []*model.DocumentChunk
	if err := list.Structs(&items); err != nil {
		return nil, err
	}

	return items, nil
}
