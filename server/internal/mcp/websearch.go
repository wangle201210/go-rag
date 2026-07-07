package mcp

import (
	"context"
	"fmt"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	v1 "github.com/wangle201210/go-rag/server/api/rag/v1"
)

type WebSearchParam struct {
	Question string `json:"question" description:"用户提问的问题，将通过 You.com 联网检索" required:"true"`
	TopK     int    `json:"top_k" description:"检索结果的数量，默认为5" required:"false"` // 默认为5
}

func GetWebSearchTool() *protocol.Tool {
	tool, err := protocol.NewTool("web_search", "You.com 联网检索，返回带来源链接的网页结果", WebSearchParam{})
	if err != nil {
		g.Log().Errorf(gctx.New(), "Failed to create tool: %v", err)
		return nil
	}
	return tool
}

func HandleWebSearch(ctx context.Context, toolReq *protocol.CallToolRequest) (res *protocol.CallToolResult, err error) {
	var req WebSearchParam
	if err := protocol.VerifyAndUnmarshal(toolReq.RawArguments, &req); err != nil {
		return nil, err
	}
	webSearch, err := c.WebSearch(ctx, &v1.WebSearchReq{
		Question: req.Question,
		TopK:     req.TopK,
	})
	if err != nil {
		// 联网检索未配置属于可预期的情况，作为友好文本返回，而不是协议错误
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				&protocol.TextContent{
					Type: "text",
					Text: fmt.Sprintf("联网检索暂不可用：%v，请配置 YDC_API_KEY 后重试", err),
				},
			},
		}, nil
	}
	docs := webSearch.Document
	msg := fmt.Sprintf("web_search %d documents", len(docs))
	for i, doc := range docs {
		msg += fmt.Sprintf("\n%d. score: %.2f, url: %v, content: %s", i+1, doc.Score(), doc.MetaData["url"], doc.Content)
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
