package mcp

import (
	"context"
	"fmt"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	v1 "github.com/wangle201210/go-rag/server/api/rag/v1"
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
	statusOK := v1.StatusOK
	getList, err := c.KBGetList(ctx, &v1.KBGetListReq{
		Status: &statusOK,
	})
	if err != nil {
		return nil, err
	}
	list := getList.List
	msg := fmt.Sprintf("get %d knowledgeBase", len(list))
	for _, l := range list {
		msg += fmt.Sprintf("\n - name: %s, description: %s", l.Name, l.Description)
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
