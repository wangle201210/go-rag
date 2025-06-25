package rag

import (
	"context"

	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/frame/g"
	v1 "github.com/wangle201210/go-rag/server/api/rag/v1"
	"github.com/wangle201210/go-rag/server/core/common"
	"github.com/wangle201210/go-rag/server/internal/logic/chat"
)

// ChatStream 流式输出接口
func (c *ControllerV1) ChatStream(ctx context.Context, req *v1.ChatStreamReq) (res *v1.ChatStreamRes, err error) {
	var streamReader *schema.StreamReader[*schema.Message]
	// 获取检索结果
	retriever, err := c.Retriever(ctx, &v1.RetrieverReq{
		Question:      req.Question,
		TopK:          req.TopK,
		Score:         req.Score,
		KnowledgeName: req.KnowledgeName,
	})
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	// 获取Chat实例
	chatI := chat.GetChat()
	// 获取流式响应
	streamReader, err = chatI.GetAnswerStream(ctx, req.ConvID, retriever.Document, req.Question)
	if err != nil {
		g.Log().Error(ctx, err)
		return &v1.ChatStreamRes{}, nil
	}
	defer streamReader.Close()
	err = common.SteamResponse(ctx, streamReader, retriever.Document)
	if err != nil {
		g.Log().Error(ctx, err)
		return
	}
	return &v1.ChatStreamRes{}, nil
}
