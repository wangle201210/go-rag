package rag

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	gorag "github.com/wangle201210/go-rag/server/core"
	"github.com/wangle201210/go-rag/server/internal/logic/knowledge"
	"github.com/wangle201210/go-rag/server/internal/logic/rag"
	"github.com/wangle201210/go-rag/server/internal/model/entity"

	v1 "github.com/wangle201210/go-rag/server/api/rag/v1"
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

	documents := entity.KnowledgeDocuments{
		KnowledgeBaseName: req.KnowledgeName,
		FileName:          req.File.Filename,
		Status:            int(v1.StatusPending),
	}
	documentsId, err := knowledge.SaveDocumentsInfo(ctx, documents)
	if err != nil {
		g.Log().Errorf(ctx, "SaveDocumentsInfo failed, err=%v", err)
		return
	}

	indexReq := &gorag.IndexReq{
		URI:           uri,
		KnowledgeName: req.KnowledgeName,
		DocumentsId:   documentsId,
	}
	ids, err := svr.Index(ctx, indexReq)
	if err != nil {
		return
	}
	res = &v1.IndexerRes{
		DocIDs: ids,
	}
	return
}
