package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/wangle201210/go-rag/server/internal/model"
)

type IKnowledgeBase interface {
	GetList(ctx context.Context, req *model.KnowledgeBaseListReq) (*model.KnowledgeBaseListRes, error)
	Create(ctx context.Context, req *model.KnowledgeBaseCreateReq) (*model.KnowledgeBaseCreateRes, error)
	Update(ctx context.Context, req *model.KnowledgeBaseUpdateReq) (*model.KnowledgeBaseUpdateRes, error)
	Delete(ctx context.Context, req *model.KnowledgeBaseDeleteReq) (*model.KnowledgeBaseDeleteRes, error)
}

type knowledgeBaseImpl struct{}

func NewKnowledgeBase() IKnowledgeBase {
	return &knowledgeBaseImpl{}
}

func (s *knowledgeBaseImpl) GetList(ctx context.Context, req *model.KnowledgeBaseListReq) (*model.KnowledgeBaseListRes, error) {
	m := g.DB().Model("knowledge_base")

	if req.Category != "" {
		m = m.Where("category", req.Category)
	}
	if req.Keyword != "" {
		m = m.WhereLike("name", "%"+req.Keyword+"%")
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	list, err := m.Page(req.Page, req.PageSize).OrderDesc("id").All()
	if err != nil {
		return nil, err
	}

	var items []model.KnowledgeBaseItem
	if err := list.Structs(&items); err != nil {
		return nil, err
	}

	return &model.KnowledgeBaseListRes{
		List:  items,
		Total: total,
	}, nil
}

func (s *knowledgeBaseImpl) Create(ctx context.Context, req *model.KnowledgeBaseCreateReq) (*model.KnowledgeBaseCreateRes, error) {
	result, err := g.DB().Model("knowledge_base").Insert(g.Map{
		"name":        req.Name,
		"description": req.Description,
		"category":    req.Category,
		"status":      1,
	})
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &model.KnowledgeBaseCreateRes{
		Id: id,
	}, nil
}

func (s *knowledgeBaseImpl) Update(ctx context.Context, req *model.KnowledgeBaseUpdateReq) (*model.KnowledgeBaseUpdateRes, error) {
	_, err := g.DB().Model("knowledge_base").
		Where("id", req.Id).
		Update(g.Map{
			"name":        req.Name,
			"description": req.Description,
			"category":    req.Category,
			"status":      req.Status,
		})
	if err != nil {
		return nil, err
	}

	return &model.KnowledgeBaseUpdateRes{
		Success: true,
	}, nil
}

func (s *knowledgeBaseImpl) Delete(ctx context.Context, req *model.KnowledgeBaseDeleteReq) (*model.KnowledgeBaseDeleteRes, error) {
	_, err := g.DB().Model("knowledge_base").Where("id", req.Id).Delete()
	if err != nil {
		return nil, err
	}

	return &model.KnowledgeBaseDeleteRes{
		Success: true,
	}, nil
}
