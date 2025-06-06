package mcp

import (
	"context"
	"fmt"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/wangle201210/go-rag/server/internal/logic/rag"
)

type KnowledgeBaseParam struct {
}

func GetKnowledgeBaseTool() *protocol.Tool {
	tool, err := protocol.NewTool("getKnowledgeBaseList", "获取知识库列表", KnowledgeBaseParam{})
	if err != nil {
		g.Log().Errorf(gctx.New(), "Failed to create tool: %v", err)
		return nil
	}
	return tool
}

func HandleKnowledgeBase(ctx context.Context, toolReq *protocol.CallToolRequest) (res *protocol.CallToolResult, err error) {
	svr := rag.GetRagSvr()
	list, err := svr.GetKnowledgeBaseList(ctx)
	msg := fmt.Sprintf("get %d knowledgeBase", len(list))
	for i, name := range list {
		msg += fmt.Sprintf("\n%d. %s", i+1, name)
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
