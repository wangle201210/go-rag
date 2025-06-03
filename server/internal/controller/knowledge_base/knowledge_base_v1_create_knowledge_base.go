package knowledge_base

import (
	"context"

	v1 "github.com/wangle201210/go-rag/server/api/knowledge_base/v1"
	"github.com/wangle201210/go-rag/server/core/model"
)

func (c *ControllerV1) CreateKnowledgeBase(ctx context.Context, req *v1.CreateKnowledgeBaseReq) (res *v1.CreateKnowledgeBaseRes, err error) {
	create, err := c.knowledgeBaseService.Create(ctx, &model.KnowledgeBaseCreateReq{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
	})
	if err != nil {
		return
	}

	res = &v1.CreateKnowledgeBaseRes{
		Id: create.Id,
	}
	return
}
