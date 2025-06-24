package rag

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	v1 "github.com/wangle201210/go-rag/server/api/rag/v1"
	"github.com/wangle201210/go-rag/server/internal/logic/knowledge"
	"github.com/wangle201210/go-rag/server/internal/logic/rag"
)

func (c *ControllerV1) DocumentsDelete(ctx context.Context, req *v1.DocumentsDeleteReq) (res *v1.DocumentsDeleteRes, err error) {
	svr := rag.GetRagSvr()

	ChunksList, err := knowledge.GetAllChunksByDocId(ctx, req.DocumentId, "id", "chunk_id")
	if err != nil {
		g.Log().Errorf(ctx, "DeleteDocumentAndChunks: GetAllChunksByDocId failed for id %d, err: %v", req.DocumentId, err)
		return
	}

	if len(ChunksList) > 0 {
		for _, chunk := range ChunksList {
			if chunk.ChunkId != "" {
				err = svr.DeleteDocument(ctx, chunk.ChunkId)
				if err != nil {
					g.Log().Errorf(ctx, "DeleteDocumentAndChunks: ES DeleteByQuery failed for docId %v, err: %v", chunk.ChunkId, err)
					return
				}
			}
		}
	}

	err = knowledge.DeleteDocument(ctx, req.DocumentId)
	return
}
