package rag

import (
	"context"

	v1 "github.com/wangle201210/go-rag/server/api/rag/v1"
	"github.com/wangle201210/go-rag/server/internal/logic/knowledge"
	"github.com/wangle201210/go-rag/server/internal/model/entity"
)

func (c *ControllerV1) UpdateChunk(ctx context.Context, req *v1.UpdateChunkReq) (res *v1.UpdateChunkRes, err error) {
	err = knowledge.UpdateChunkByIds(ctx, req.Ids, entity.KnowledgeChunks{
		Status: req.Status,
	})
	if err != nil {
		return
	}

	return
}
