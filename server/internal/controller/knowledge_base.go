package controller

import (
	"context"

	"github.com/wangle201210/go-rag/server/core/model"
	"github.com/wangle201210/go-rag/server/core/server"
)

type KnowledgeBaseController struct {
	knowledgeBaseService server.IKnowledgeBase
}

func NewKnowledgeBaseController() *KnowledgeBaseController {
	return &KnowledgeBaseController{
		knowledgeBaseService: server.NewKnowledgeBase(),
	}
}

func (c *KnowledgeBaseController) GetList(ctx context.Context, req *model.KnowledgeBaseListReq) (*model.KnowledgeBaseListRes, error) {
	return c.knowledgeBaseService.GetList(ctx, req)
}

func (c *KnowledgeBaseController) Create(ctx context.Context, req *model.KnowledgeBaseCreateReq) (*model.KnowledgeBaseCreateRes, error) {
	return c.knowledgeBaseService.Create(ctx, req)
}

func (c *KnowledgeBaseController) Update(ctx context.Context, req *model.KnowledgeBaseUpdateReq) (*model.KnowledgeBaseUpdateRes, error) {
	return c.knowledgeBaseService.Update(ctx, req)
}

func (c *KnowledgeBaseController) Delete(ctx context.Context, req *model.KnowledgeBaseDeleteReq) (*model.KnowledgeBaseDeleteRes, error) {
	return c.knowledgeBaseService.Delete(ctx, req)
}
