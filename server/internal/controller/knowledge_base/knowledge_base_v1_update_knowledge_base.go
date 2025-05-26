package knowledge_base

import (
	"context"

	"github.com/wangle201210/go-rag/server/api/knowledge_base/v1"
	"github.com/wangle201210/go-rag/server/internal/model"
)

func (c *ControllerV1) UpdateKnowledgeBase(ctx context.Context, req *v1.UpdateKnowledgeBaseReq) (res *v1.UpdateKnowledgeBaseRes, err error) {
	update, err := c.knowledgeBaseService.Update(ctx, &model.KnowledgeBaseUpdateReq{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Status:      req.Status,
	})
	if err != nil {
		return
	}

	res = &v1.UpdateKnowledgeBaseRes{
		Success: update.Success,
	}
	return
}
