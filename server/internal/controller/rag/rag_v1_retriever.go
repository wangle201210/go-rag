package rag

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/gogf/gf/v2/frame/g"
	gorag "github.com/wangle201210/go-rag/server/core"
	"github.com/wangle201210/go-rag/server/internal/logic/rag"

	"github.com/wangle201210/go-rag/server/api/rag/v1"
)

func (c *ControllerV1) Retriever(ctx context.Context, req *v1.RetrieverReq) (res *v1.RetrieverRes, err error) {
	ragSvr := rag.GetRagSvr()
	if req.TopK == 0 {
		req.TopK = 5
	}
	if req.Score == 0 {
		req.Score = 0.2
	}
	if req.Score < 1.0 {
		req.Score += 1
	}
	ragReq := &gorag.RetrieveReq{
		Query:         req.Question,
		TopK:          req.TopK,
		Score:         req.Score,
		KnowledgeName: req.KnowledgeName,
	}
	g.Log().Infof(ctx, "ragReq: %v", ragReq)
	msg, err := ragSvr.Retrieve(ctx, ragReq)
	if err != nil {
		return
	}
	for _, document := range msg {
		if document.MetaData != nil {
			delete(document.MetaData, "_dense_vector")
			m := make(map[string]interface{})
			if err = json.Unmarshal([]byte(document.MetaData["ext"].(string)), &m); err != nil {
				return
			}
			document.MetaData["ext"] = m
		}
	}
	// eino 默认是把分高的排在两边，这里我xiu gai
	sort.Slice(msg, func(i, j int) bool {
		return msg[i].Score() > msg[j].Score()
	})
	res = &v1.RetrieverRes{
		Document: msg,
	}
	return
}
