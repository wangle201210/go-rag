package rag

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/wangle201210/go-rag/server/api/rag/v1"
)

func (c *ControllerV1) RetrieverDify(ctx context.Context, req *v1.RetrieverDifyReq) (res *v1.RetrieverDifyRes, err error) {
	retriever, err := c.Retriever(ctx, &v1.RetrieverReq{
		Question:      req.Query,
		TopK:          req.RetrievalSetting.TopK,
		Score:         req.RetrievalSetting.ScoreThreshold,
		KnowledgeName: req.KnowledgeID,
	})
	if err != nil {
		return
	}
	res = &v1.RetrieverDifyRes{}
	for _, document := range retriever.Document {
		g.Log().Infof(ctx, "content: %s, score: %f", document.Content, document.Score())
		record := &v1.Record{
			Content: document.Content,
			Score:   document.Score(),
		}
		res.Records = append(res.Records, record)
	}
	return
}
