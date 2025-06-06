package mcp

import (
	"context"
	"fmt"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	v1 "github.com/wangle201210/go-rag/server/api/rag/v1"
)

type RetrieverParam struct {
	Question      string  `json:"question" description:"用户提问的问题" required:"true"`
	KnowledgeName string  `json:"knowledge_name" description:"知识库名称，请先通过getKnowledgeBaseList获取列表后判断是否有符合用户提示词的知识库" required:"true"`
	TopK          int     `json:"top_k" description:"检索结果的数量，默认为10" required:"false"`     // 默认为5
	Score         float64 `json:"score"  description:"检索结果的分数阀值，默认为0.2" required:"false"` // 默认为0.2
}

func GetRetrieverTool() *protocol.Tool {
	tool, err := protocol.NewTool("retriever", "检索知识库文档", RetrieverParam{})
	if err != nil {
		g.Log().Errorf(gctx.New(), "Failed to create tool: %v", err)
		return nil
	}
	return tool
}

func HandleRetriever(ctx context.Context, toolReq *protocol.CallToolRequest) (res *protocol.CallToolResult, err error) {
	var req RetrieverParam
	if err := protocol.VerifyAndUnmarshal(toolReq.RawArguments, &req); err != nil {
		return nil, err
	}
	retriever, err := c.Retriever(ctx, &v1.RetrieverReq{
		Question:      req.Question,
		TopK:          req.TopK,
		Score:         req.Score,
		KnowledgeName: req.KnowledgeName,
	})
	if err != nil {
		return nil, err
	}
	docs := retriever.Document
	msg := fmt.Sprintf("retrieve %d documents", len(docs))
	for i, doc := range docs {
		msg += fmt.Sprintf("\n%d. score: %.2f, content: %s", i+1, doc.Score(), doc.Content)
	}
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: msg,
			},
		},
	}, nil
}
