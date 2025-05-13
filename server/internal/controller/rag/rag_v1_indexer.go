package rag

import (
	"context"

	"github.com/wangle201210/go-rag/server/internal/logic/rag"

	"github.com/wangle201210/go-rag/server/api/rag/v1"
)

func (c *ControllerV1) Indexer(ctx context.Context, req *v1.IndexerReq) (res *v1.IndexerRes, err error) {
	svr := rag.GetRagSvr()
	uri := req.URL
	if req.File != nil {
		filename, e := req.File.Save("./uploads/")
		if e != nil {
			err = e
			return
		}
		uri = "./uploads/" + filename
	}
	ids, err := svr.Index(uri)
	if err != nil {
		return
	}
	res = &v1.IndexerRes{
		DocIDs: ids,
	}
	return
}
