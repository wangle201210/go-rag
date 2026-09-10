package rag

import (
	"context"
	"sort"

	"github.com/cloudwego/eino/components/retriever"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/wangle201210/go-rag/server/internal/logic/websearch"

	v1 "github.com/wangle201210/go-rag/server/api/rag/v1"
)

// WebSearch 通过 You.com Search API 实现联网检索，独立于知识库向量检索
func (c *ControllerV1) WebSearch(ctx context.Context, req *v1.WebSearchReq) (res *v1.WebSearchRes, err error) {
	svr := websearch.GetWebSearchSvr()
	if svr == nil {
		return nil, gerror.New("联网检索未配置，请设置环境变量 YDC_API_KEY 或在配置文件 websearch.youcom.apiKey 中填写 You.com API Key")
	}
	g.Log().Infof(ctx, "webSearchReq: %v", req)
	docs, err := svr.Retrieve(ctx, req.Question, retriever.WithTopK(req.TopK))
	if err != nil {
		return nil, err
	}
	sort.Slice(docs, func(i, j int) bool {
		return docs[i].Score() > docs[j].Score()
	})
	res = &v1.WebSearchRes{
		Document: docs,
	}
	return
}
