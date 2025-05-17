package rag

import (
	"context"

	"github.com/wangle201210/go-rag/server/api/rag/v1"
	"github.com/wangle201210/go-rag/server/internal/logic/chat"
)

func (c *ControllerV1) Chat(ctx context.Context, req *v1.ChatReq) (res *v1.ChatRes, err error) {
	retriever, err := c.Retriever(ctx, &v1.RetrieverReq{
		Question:      req.Question,
		TopK:          req.TopK,
		Score:         req.Score,
		KnowledgeName: req.KnowledgeName,
	})
	if err != nil {
		return
	}
	chatI := chat.GetChat()
	answer, err := chatI.GetAnswer(ctx, req.ConvID, retriever.Document, req.Question)
	if err != nil {
		return
	}
	res = &v1.ChatRes{
		Answer:     answer,
		References: retriever.Document,
	}
	return
}
