package retriever

import (
	"context"

	"github.com/wangle201210/go-rag/server/internal/logic"

	"github.com/wangle201210/go-rag/server/api/retriever/v1"
)

func (c *ControllerV1) Retriever(ctx context.Context, req *v1.RetrieverReq) (res *v1.RetrieverRes, err error) {
	ragSvr := logic.GetRagSvr()
	if req.Score < 1.0 {
		req.Score += 1
	}
	msg, err := ragSvr.Retrieve(req.Question, req.Score, req.TopK)
	if err != nil {
		return
	}
	for _, document := range msg {
		if document.MetaData != nil {
			delete(document.MetaData, "_dense_vector")
			// if v, e := document.MetaData["_score"]; e {
			// 	vf := v.(float64)
			// 	document.MetaData["_score"] = vf - 1
			// }
		}
	}
	res = &v1.RetrieverRes{
		Document: msg,
	}
	return
}
