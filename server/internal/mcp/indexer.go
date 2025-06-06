package mcp

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	gorag "github.com/wangle201210/go-rag/server/core"
	"github.com/wangle201210/go-rag/server/internal/logic/rag"
)

type IndexParam struct {
	URI           string `json:"uri" description:"文件路径" required:"true"` // 可以是文件路径（pdf，html，md等），也可以是网址
	KnowledgeName string `json:"knowledge_name" description:"知识库名字,请先通过getKnowledgeBaseList获取列表后判断是否有符合的知识库，如果没有则根据用户提示词自己生成" required:"true"`
}

func GetIndexerByFilePathTool() *protocol.Tool {
	tool, err := protocol.NewTool("Indexer_by_filepath", "通过文件路径进行文本嵌入", IndexParam{})
	if err != nil {
		g.Log().Errorf(gctx.New(), "Failed to create tool: %v", err)
		return nil
	}
	return tool
}

func HandleIndexerByFilePath(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	var reqData IndexParam
	if err := protocol.VerifyAndUnmarshal(req.RawArguments, &reqData); err != nil {
		return nil, err
	}
	svr := rag.GetRagSvr()
	uri := reqData.URI
	indexReq := &gorag.IndexReq{
		URI:           uri,
		KnowledgeName: reqData.KnowledgeName,
	}
	ids, err := svr.Index(ctx, indexReq)
	if err != nil {
		return nil, err
	}
	msg := fmt.Sprintf("index file %s successfully, knowledge_name: %s, doc_ids: %v", uri, reqData.KnowledgeName, ids)
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: msg,
			},
		},
	}, nil
}

type IndexFileParam struct {
	Filename      string `json:"filename" description:"文件名字" required:"true"`
	Content       string `json:"content" description:"被base64编码后的文件内容，先调用工具获取base64信息" required:"true"` // 可以是文件路径（pdf，html，md等），也可以是网址文件" required:"true"` // 可以是文件路径（pdf，html，md等），也可以是网址
	KnowledgeName string `json:"knowledge_name" description:"知识库名字,请先通过getKnowledgeBaseList获取列表后判断是否有符合的知识库，如果没有则根据用户提示词自己生成" required:"true"`
}

func GetIndexerByFileBase64ContentTool() *protocol.Tool {
	tool, err := protocol.NewTool("Indexer_by_base64_file_content", "获取文件base64信息后上传，然后对内容进行文本嵌入", IndexFileParam{})
	if err != nil {
		g.Log().Errorf(gctx.New(), "Failed to create tool: %v", err)
		return nil
	}
	return tool
}

func HandleIndexerByFileBase64Content(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	var reqData IndexFileParam
	if err := protocol.VerifyAndUnmarshal(req.RawArguments, &reqData); err != nil {
		return nil, err
	}
	// svr := rag.GetRagSvr()
	decoded, err := base64.StdEncoding.DecodeString(reqData.Content)
	if err != nil {
		return nil, err
	}
	fmt.Println(decoded)
	// indexReq := &gorag.IndexReq{
	// 	URI:           uri,
	// 	KnowledgeName: reqData.KnowledgeName,
	// }
	// ids, err := svr.Index(ctx, indexReq)
	// if err != nil {
	// 	return nil, err
	// }
	// msg := fmt.Sprintf("index file %s successfully, knowledge_name: %s, doc_ids: %v", uri, reqData.KnowledgeName, ids)
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: string(decoded),
			},
		},
	}, nil
}
