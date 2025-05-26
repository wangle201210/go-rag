package knowledge_base

import (
	"context"

	"github.com/wangle201210/go-rag/server/internal/model"

	"github.com/wangle201210/go-rag/server/api/knowledge_base/v1"
)

func (c *ControllerV1) DeleteKnowledgeBase(ctx context.Context, req *v1.DeleteKnowledgeBaseReq) (res *v1.DeleteKnowledgeBaseRes, err error) {
	deleteRes, err := c.knowledgeBaseService.Delete(ctx, &model.KnowledgeBaseDeleteReq{
		Id: req.Id,
	})
	if err != nil {
		return
	}
	res = &v1.DeleteKnowledgeBaseRes{
		Success: deleteRes.Success,
	}
	return
}
