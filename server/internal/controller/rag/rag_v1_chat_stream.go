package rag

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/google/uuid"
	v1 "github.com/wangle201210/go-rag/server/api/rag/v1"
	"github.com/wangle201210/go-rag/server/core/common"
	"github.com/wangle201210/go-rag/server/internal/logic/chat"
)

// ChatStream 流式输出接口
func (c *ControllerV1) ChatStream(ctx context.Context, req *v1.ChatStreamReq) (res *v1.ChatStreamRes, err error) {
	// 获取HTTP响应对象
	httpReq := ghttp.RequestFromCtx(ctx)
	httpResp := httpReq.Response
	// 设置响应头
	httpResp.Header().Set("Content-Type", "text/event-stream")
	httpResp.Header().Set("Cache-Control", "no-cache")
	httpResp.Header().Set("Connection", "keep-alive")
	httpResp.Header().Set("X-Accel-Buffering", "no") // 禁用Nginx缓冲
	httpResp.Header().Set("Access-Control-Allow-Origin", "*")

	// 获取检索结果
	retriever, err := c.Retriever(ctx, &v1.RetrieverReq{
		Question:      req.Question,
		TopK:          req.TopK,
		Score:         req.Score,
		KnowledgeName: req.KnowledgeName,
	})
	if err != nil {
		writeSSEError(httpResp, err)
		return &v1.ChatStreamRes{}, nil
	}
	sd := &common.StreamData{
		Id:       uuid.NewString(),
		Created:  time.Now().Unix(),
		Document: retriever.Document,
	}
	marshal, _ := sonic.Marshal(sd)
	writeSSEDocuments(httpResp, string(marshal))
	// 获取Chat实例
	chatI := chat.GetChat()
	// 获取流式响应
	streamReader, err := chatI.GetAnswerStream(ctx, req.ConvID, retriever.Document, req.Question)
	if err != nil {
		writeSSEError(httpResp, err)
		return &v1.ChatStreamRes{}, nil
	}
	defer streamReader.Close()

	// 处理流式响应
	for {
		chunk, err := streamReader.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			writeSSEError(httpResp, err)
			break
		}
		if len(chunk.Content) == 0 {
			continue
		}

		sd.Content = chunk.Content
		marshal, _ := sonic.Marshal(sd)
		// 发送数据事件
		writeSSEData(httpResp, string(marshal))
	}
	// 发送结束事件
	writeSSEDone(httpResp)
	return &v1.ChatStreamRes{}, nil
}

// writeSSEData 写入SSE事件
func writeSSEData(resp *ghttp.Response, data string) {
	if len(data) == 0 {
		return
	}
	// g.Log().Infof(context.Background(), "data: %s", data)
	resp.Writeln(fmt.Sprintf("data:%s\n", data))
	resp.Flush()
}

func writeSSEDone(resp *ghttp.Response) {
	resp.Writeln(fmt.Sprintf("data:%s\n", "[DONE]"))
	resp.Flush()
}

func writeSSEDocuments(resp *ghttp.Response, data string) {
	resp.Writeln(fmt.Sprintf("documents:%s\n", data))
	resp.Flush()
}

// writeSSEError 写入SSE错误
func writeSSEError(resp *ghttp.Response, err error) {
	g.Log().Error(context.Background(), err)
	resp.Writeln(fmt.Sprintf("event: error\ndata: %s\n\n", err.Error()))
	resp.Flush()
}
